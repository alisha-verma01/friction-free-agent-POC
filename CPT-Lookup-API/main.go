package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Configuration constants
const (
	ExternalAPIURL = "https://SECRET_APIreferencdcodes/v1/policy/codes/v4"
	TokenURL       = "https://SECRET_TOKEN/auth/oauth2/token"
	TokenExpiryBufferTime = 30 * time.Second
	TokenRefreshBufferTime = 30 * time.Second
)

// Static disclaimer text used for fallback responses
const DefaultDisclaimerText = "The search executed is based on data that you have selected. Your search is not a request for prior authorization and is not notification to UnitedHealthcare. Prior authorization is required for services that require it."

var HttpClient = &http.Client{
	Timeout: 4 * time.Second,
}

// TokenResponse struct
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type TokenManager struct {
	mutex       sync.RWMutex
	token       string
	expiresAt   time.Time
	clientID    string
	clientSecret string
}

var tokenManager *TokenManager

func NewTokenManager() *TokenManager {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("CLIENT_ID and CLIENT_SECRET environment variables must be set")
	}

	tm := &TokenManager{
		clientID:    clientID,
		clientSecret: clientSecret,
	}

	// Get initial token
	if err := tm.RefreshToken(); err != nil {
		log.Fatalf("Failed to get initial token: %v", err)
	}

	return tm
}

func (tm *TokenManager) RefreshToken() error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", tm.clientID)
	data.Set("client_secret", tm.clientSecret)

	req, err := http.NewRequest("POST", TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("token request failed: " + string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	tm.token = tokenResp.AccessToken
	tm.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	log.Printf("Token refreshed successfully, expires at: %v", tm.expiresAt)
	return nil
}

func (tm *TokenManager) GetValidToken() (string, error) {
	tm.mutex.RLock()
	if time.Until(tm.expiresAt) <= TokenRefreshBufferTime {
		tm.mutex.RUnlock()
		if err := tm.RefreshToken(); err != nil {
			return "", err
		}
		tm.mutex.RLock()
	}
	token := tm.token
	tm.mutex.RUnlock()
	return "Bearer " + token, nil
}

type ProxyRequest struct {
	PolicyIssueState string   `json:"policyIssueState"`
	ProcedureCode    []string `json:"procedureCode"`
}

type ExternalRequest struct {
	ProductCategory   string   `json:"productCategory"`
	PolicyIssueState  string   `json:"policyIssueState"`
	ProcedureCode     []string `json:"procedureCode"`
	TIN               string   `json:"tin"`
}

func storeAPIResponse(db *gorm.DB, apiData models.APIResponse) error {
	var existingResponse models.APIResponse
	err := db.Where("product_category = ? AND policy_issue_state = ?", 
		apiData.ProductCategory, apiData.PolicyIssueState).
		Preload("PreliminaryDeterminations.SiteOfServices.Conditions").
		First(&existingResponse).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return db.Create(&apiData).Error
	} else if err != nil {
		return err
	}

	if time.Since(existingResponse.UpdatedAt) > DataExpiryHours*time.Hour {
		db.Where("api_response_id = ?", existingResponse.ID).
			Delete(&models.PreliminaryDetermination{})
		existingResponse.ProductCategory = apiData.ProductCategory
		existingResponse.PolicyIssueState = apiData.PolicyIssueState
		existingResponse.DiagnosisCode = apiData.DiagnosisCode
		existingResponse.UpdatedAt = time.Now()
	}

	for _, pd := range apiData.PreliminaryDeterminations {
		pd.ApiResponseID = existingResponse.ID
		existingResponse.PreliminaryDeterminations = append(existingResponse.PreliminaryDeterminations, pd)
	}

	return db.Save(&existingResponse).Error
}

func getLocalData(db *gorm.DB, policyIssueState string, procedureCodes []string) (models.APIResponse, error) {
	var response models.APIResponse

	result := db.Where("policy_issue_state = ?", policyIssueState).
		Preload("PreliminaryDeterminations", func(db *gorm.DB) *gorm.DB {
			if len(procedureCodes) > 0 {
				return db.Where("procedure_code IN ?", procedureCodes)
			}
			return db
		}).
		Preload("PreliminaryDeterminations.SiteOfServices.Conditions").
		First(&response)

	if result.Error != nil {
		return response, result.Error
	}

	if len(procedureCodes) > 0 && len(response.PreliminaryDeterminations) == 0 {
		return response, gorm.ErrRecordNotFound
	}

	return response, nil
}

func proxyHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody ProxyRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Received request: PolicyIssueState=%s, ProcedureCodes=%v",
			reqBody.PolicyIssueState, reqBody.ProcedureCode)

		extReq := ExternalRequest{
			ProductCategory:  ProductCategory,
			PolicyIssueState: reqBody.PolicyIssueState,
			ProcedureCode:    reqBody.ProcedureCode,
			TIN:              TIN,
		}

		jsonData, err := json.Marshal(extReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		req, err := http.NewRequest("POST", ExternalAPIURL, strings.NewReader(string(jsonData)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		xIdentity := os.Getenv("X_IDENTITY")
		if xIdentity == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "X_IDENTITY environment variable not set"})
			return
		}

		authToken, err := tokenManager.GetValidToken()
		if err != nil {
			log.Printf("Failed to get valid token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get authentication token"})
			return
		}

		req.Header.Set("Authorization", authToken)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("X-Identity", xIdentity)
		req.Header.Set("Content-Type", "application/json")

		resp, err := HttpClient.Do(req)
		if err != nil {
			log.Printf("External API call failed: %v, attempting to use local data", err)
			if localData := tryLocalFallback(db, reqBody, c); localData {
				return
			}
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":  "External API unavailable and no local data found",
				"detail": err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		log.Printf("External API response status: %d", resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read API response: %v, attempting to use local data", err)
			if localData := tryLocalFallback(db, reqBody, c); localData {
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read API response and no local data available"})
			return
		}

		if resp.StatusCode >= 400 {
			log.Printf("External API returned error status %d, response body: %s", resp.StatusCode, string(body))
			if localData := tryLocalFallback(db, reqBody, c); localData {
				return
			}
			c.Data(resp.StatusCode, "application/json", body)
			return
		}

		if resp.StatusCode == http.StatusOK {
			var apiResponseData models.APIResponse
			if err := json.Unmarshal(body, &apiResponseData); err != nil {
				log.Printf("Failed to parse API response for database storage: %v", err)
				c.Data(resp.StatusOK, "application/json", body)
				return
			}

			if err := storeAPIResponse(db, apiResponseData); err != nil {
				log.Printf("Failed to store API response in database: %v", err)
			}

			c.JSON(http.StatusOK, apiResponseData)
			return
		}

		c.Data(resp.StatusCode, "application/json", body)
	}
}
func tryLocalFallback(db *gorm.DB, reqBody ProxyRequest, c *gin.Context) bool {
	localData, dbErr := getLocalData(db, reqBody.PolicyIssueState, reqBody.ProcedureCode)
	if dbErr != nil {
		log.Printf("No local data available: %v", dbErr)
		return false
	}

	addGoldCardStatus(db, localData)

	localData.DisclaimerText = DefaultDisclaimerText
	c.JSON(http.StatusOK, localData)
	return true
}

func addGoldCardStatus(db *gorm.DB, apiData models.APIResponse) {
	cptCodes := make([]string, 0, len(apiData.PreliminaryDeterminations))
	for _, pd := range apiData.PreliminaryDeterminations {
		cptCodes = append(cptCodes, pd.ProcedureCode)
	}

	var goldCardCodes []models.GoldCardCode
	db.Where("cpt_code IN ?", cptCodes).Find(&goldCardCodes)

	goldCardMap := make(map[string]bool)
	for _, gcc := range goldCardCodes {
		goldCardMap[gcc.CptCode] = true
	}

	for i := range apiData.PreliminaryDeterminations {
		apiData.PreliminaryDeterminations[i].IsGoldCard = goldCardMap[apiData.PreliminaryDeterminations[i].ProcedureCode]
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	tokenManager = NewTokenManager()
	router := gin.Default()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(
		&models.APIResponse{},
		&models.PreliminaryDetermination{},
		&models.SiteOfService{},
		&models.Condition{},
		&models.GoldCardCode{},
	)
	if err != nil {
		panic(err)
	}

	router.POST("/api", proxyHandler(db))

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}

	err = router.Run(port)
	if err != nil {
		panic(err)
	}
}
