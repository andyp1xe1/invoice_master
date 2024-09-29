package main

import (
	"gorm.io/gorm"
	"time"
)

// Service model
type Service struct {
	gorm.Model
	InvoiceID  *uint   `gorm:"column:invoice_id;index"`  // Optional reference to Invoice
	ContractID *uint   `gorm:"column:contract_id;index"` // Optional reference to Contract
	Item       string  `gorm:"column:item;size:255;not null"`
	Quantity   int     `gorm:"column:quantity;not null"`
	PriceUnit  float64 `gorm:"column:price_unit;not null"`
	Taxes      float64 `gorm:"column:taxes"`
	Amount     float64 `gorm:"column:amount;not null"`
	Frequency  int     `gorm:"column:frequency;not null"`
}

type ServiceDTO struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	InvoiceID uint    `gorm:"column:invoice_id" json:"invoiceId"` // Foreign key to Invoice
}

// Contract model
type Contract struct {
	gorm.Model
	ID                      uint      `gorm:"primaryKey" json:"id"` // Primary key for Invoice
	URLHash                 string    `json:"urlHash" gorm:"column:url_hash"`
	DocID                   string    `json:"docId" gorm:"column:doc_id"`
	CurrencySymbol          string    `json:"currencySymbol" gorm:"column:currency_symbol"`
	LanguageCode            string    `json:"languageCode" gorm:"column:language_code"`
	IssueDate               time.Time `json:"issueDate" gorm:"column:issue_date"`
	DueDate                 time.Time `json:"dueDate" gorm:"column:due_date"`
	PONumber                string    `json:"PONumber" gorm:"column:po_number"`
	FromName                string    `json:"fromName" gorm:"column:from_name"`
	FromEmail               string    `json:"fromEmail" gorm:"column:from_email"`
	ToName                  string    `json:"toName" gorm:"column:to_name"`
	ToPhone                 string    `json:"toPhone" gorm:"column:to_phone"`
	ToAddress               string    `json:"toAddress" gorm:"column:to_address"`
	AdditionalInfo          string    `json:"additionalInfo" gorm:"column:additional_info"`
	CompanyPromoInfoPhone   string    `json:"companyPromoInfoPhone" gorm:"column:company_promo_info_phone"`
	CompanyPromoInfoEmail   string    `json:"companyPromoInfoEmail" gorm:"column:company_promo_info_email"`
	CompanyPromoInfoWebPage string    `json:"companyPromoInfoWebPage" gorm:"column:company_promo_info_web_page"`
	Tax                     float64   `json:"tax" gorm:"column:tax"`
	Discount                float64   `json:"discount" gorm:"column:discount"`
	Total                   float64   `json:"total" gorm:"column:total"`
	Subtotal                float64   `json:"subtotal" gorm:"column:subtotal"`
	BalanceDue              float64   `json:"balanceDue" gorm:"column:balance_due"`
	PaidDate                time.Time `json:"paidDate" gorm:"column:paid_date"`
	PaidAmount              float64   `json:"paidAmount" gorm:"column:paid_amount"`
	IsExpense               bool      `json:"isExpense" gorm:"column:is_expense"`
	ToIban                  string    `json:"toIban" gorm:"column:to_iban"`
	ToBank                  string    `json:"toBank" gorm:"column:to_bank"`
	ToBankCode              string    `json:"toBankCode" gorm:"column:to_bank_code"`
	ToIdno                  string    `json:"toIdno" gorm:"column:to_idno"`
	ToTva                   string    `json:"toTva" gorm:"column:to_tva"`
	Taxes                   float64   `json:"taxes" gorm:"column:taxes"`
	Services                []Service `json:"services" gorm:"foreignKey:InvoiceID"` // Relationship with foreign key
}
