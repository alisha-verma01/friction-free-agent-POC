package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"cpt-api/cpt-models/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&models.GoldCardCode{})
	if err != nil {
		log.Fatalf("failed to migrate GoldCardCode table: %v", err)
	}

	err = importCPTCodes(db, "gold_card/UHG-Goldcard-CPT-Codes.csv")
	if err != nil {
		log.Fatalf("failed to import CPT codes: %v", err)
	}

	fmt.Println("Successfully imported gold card CPT codes!")
}

func importCPTCodes(db *gorm.DB, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var cptCodes []models.GoldCardCode
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		cptCode := models.GoldCardCode{
			Code: line,
		}
		cptCodes = append(cptCodes, cptCode)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	if len(cptCodes) > 0 {
		result := db.CreateInBatches(cptCodes, 100)
		if result.Error != nil {
			return fmt.Errorf("failed to insert CPT codes: %v", result.Error)
		}

		fmt.Printf("Processed %d CPT codes from file\n", len(cptCodes))
		fmt.Printf("Inserted/Updated %d records in database\n", result.RowsAffected)
	}

	return nil
}
