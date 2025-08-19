package models

import (
	"time"

	"gorm.io/gorm"
)

type ApiResponse struct {
	ID                       uint                      `gorm:"primaryKey" json:"-"`
	ProductCategory          string                    `json:"productCategory"`
	PolicyIssueState         string                    `json:"policyIssueState"`
	DiagnosisCode            *string                   `json:"diagnosisCode,omitempty"`
	DisclaimerText           string                    `gorm:"-" json:"disclaimerText"`
	PreliminaryDeterminations []PreliminaryDetermination `gorm:"foreignKey:ApiResponseID" json:"preliminaryDeterminations"`
	CreatedAt                time.Time                 `json:"-"`
	UpdatedAt                time.Time                 `json:"-"`
}

func (a *ApiResponse) BeforeCreate(tx *gorm.DB) (err error) {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return
}

func (a *ApiResponse) BeforeUpdate(tx *gorm.DB) (err error) {
	a.UpdatedAt = time.Now()
	return
}

type PreliminaryDetermination struct {
	ID                   uint             `gorm:"primaryKey" json:"-"`
	ApiResponseID        uint             `json:"-"`
	ProcedureCode        string           `json:"procedureCode"`
	ProcedureDescription string           `json:"procedureCodeDesc"`
	DecisionSummaryText  string           `json:"decisionSummaryText"`
	IsGoldCard           bool             `gorm:"-" json:"gold_card_code"`
	SiteOfServices       []SiteOfService  `gorm:"foreignKey:PreliminaryDeterminationID" json:"siteOfServices"`
	CreatedAt            time.Time        `json:"-"`
	UpdatedAt            time.Time        `json:"-"`
}

func (pd *PreliminaryDetermination) BeforeCreate(tx *gorm.DB) (err error) {
	pd.CreatedAt = time.Now()
	pd.UpdatedAt = time.Now()
	return
}

func (pd *PreliminaryDetermination) BeforeUpdate(tx *gorm.DB) (err error) {
	pd.UpdatedAt = time.Now()
	return
}

type SiteOfService struct {
	ID                       uint        `gorm:"primaryKey" json:"-"`
	PreliminaryDeterminationID uint     `json:"-"`
	SiteOfServiceType        string      `json:"siteOfServiceType"`
	DecisionCode             string      `json:"decisionCode"`
	Conditions               []Condition `gorm:"foreignKey:SiteOfServiceID" json:"conditions"`
	CreatedAt                time.Time   `json:"-"`
	UpdatedAt                time.Time   `json:"-"`
}

func (sos *SiteOfService) BeforeCreate(tx *gorm.DB) (err error) {
	sos.CreatedAt = time.Now()
	sos.UpdatedAt = time.Now()
	return
}

func (sos *SiteOfService) BeforeUpdate(tx *gorm.DB) (err error) {
	sos.UpdatedAt = time.Now()
	return
}

type Condition struct {
	ID              uint      `gorm:"primaryKey" json:"-"`
	SiteOfServiceID uint      `json:"-"`
	ConditionName   string    `json:"conditionName"`
	ConditionDetail string    `json:"conditionDetail"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}

func (c *Condition) BeforeCreate(tx *gorm.DB) (err error) {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return
}

func (c *Condition) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return
}

type GoldCardCode struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CPTCode   string    `gorm:"uniqueIndex;not null" json:"cpt_code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (gcc *GoldCardCode) BeforeCreate(tx *gorm.DB) (err error) {
	gcc.CreatedAt = time.Now()
	gcc.UpdatedAt = time.Now()
	return
}

func (gcc *GoldCardCode) BeforeUpdate(tx *gorm.DB) (err error) {
	gcc.UpdatedAt = time.Now()
	return
}
