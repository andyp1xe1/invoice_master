package models

import (
	"gorm.io/gorm"
	"time"
)

// Invoice model
type Invoice struct {
	gorm.Model
	UserID                  int       `gorm:"not null"`
	URLHash                 string    `gorm:"size:255;not null"`
	DocID                   string    `gorm:"size:255;not null"`
	CurrencySymbol          string    `gorm:"size:10;not null"`
	LanguageCode            string    `gorm:"size:10"`
	IssueDate               time.Time `gorm:"not null"`
	DueDate                 time.Time `gorm:"not null"`
	PONumber                string    `gorm:"size:50"`
	FromName                string    `gorm:"size:255"`
	FromEmail               string    `gorm:"size:100"`
	ToName                  string    `gorm:"size:255"`
	ToEmail                 string    `gorm:"size:100"`
	ToPhone                 string    `gorm:"size:20"`
	ToAddress               string    `gorm:"size:255"`
	Services                []Service `gorm:"foreignKey:InvoiceID"` // One-to-many relationship
	BankDetails             string    `gorm:"size:255"`
	AdditionalInfo          string    `gorm:"type:text"`
	CompanyPromoInfoPhone   string    `gorm:"size:50"`
	CompanyPromoInfoEmail   string    `gorm:"size:100"`
	CompanyPromoInfoWebPage string    `gorm:"size:255"`
	Tax                     float64
	Discount                float64
	Total                   float64 `gorm:"not null"`
	Subtotal                float64 `gorm:"not null"`
	BalanceDue              float64
	PaidDate                time.Time
	PaidAmount              float64
	IsExpense               bool
}

// Service model (used for both Invoice and Contract)
type Service struct {
	gorm.Model
	InvoiceID  *uint   `gorm:"index"` // Optional reference to Invoice
	ContractID *uint   `gorm:"index"` // Optional reference to Contract
	Item       string  `gorm:"size:255;not null"`
	Quantity   int     `gorm:"not null"`
	PriceUnit  float64 `gorm:"not null"`
	Taxes      float64
	Amount     float64 `gorm:"not null"`
}

// Contract model
type Contract struct {
	gorm.Model
	UserID                  int       `gorm:"not null"`
	URLHash                 string    `gorm:"size:255;not null"`
	DocID                   string    `gorm:"size:255;not null"`
	CurrencySymbol          string    `gorm:"size:10;not null"`
	LanguageCode            string    `gorm:"size:10"`
	IssueDate               time.Time `gorm:"not null"`
	DueDate                 time.Time `gorm:"not null"`
	PONumber                string    `gorm:"size:50"`
	FromName                string    `gorm:"size:255"`
	FromEmail               string    `gorm:"size:100"`
	ToName                  string    `gorm:"size:255"`
	ToEmail                 string    `gorm:"size:100"`
	ToPhone                 string    `gorm:"size:20"`
	ToAddress               string    `gorm:"size:255"`
	Services                []Service `gorm:"foreignKey:ContractID"` // One-to-many relationship
	BankDetails             string    `gorm:"size:255"`
	AdditionalInfo          string    `gorm:"type:text"`
	CompanyPromoInfoPhone   string    `gorm:"size:50"`
	CompanyPromoInfoEmail   string    `gorm:"size:100"`
	CompanyPromoInfoWebPage string    `gorm:"size:255"`
	Frequency               int       `gorm:"not null"`
	Subtotal                float64   `gorm:"not null"`
	Taxes                   float64
	Total                   float64 `gorm:"not null"`
}
