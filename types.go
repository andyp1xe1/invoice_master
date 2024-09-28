package main

type InvoiceDTO struct {
	ID                      int        `json:"id"`
	UserID                  int        `json:"userId"`
	Logo                    string     `json:"logo"`
	URLHash                 string     `json:"urlHash"`
	DocID                   string     `json:"docId"`
	CurrencySymbol          string     `json:"currencySymbol"`
	LanguageCode            string     `json:"languageCode"`
	IssueDate               string     `json:"issueDate"`
	DueDate                 string     `json:"dueDate"`
	PONumber                string     `json:"PONumber"`
	FromName                string     `json:"fromName"`
	FromEmail               string     `json:"fromEmail"`
	ToName                  string     `json:"toName"`
	ToEmail                 string     `json:"toEmail"`
	ToPhone                 string     `json:"toPhone"`
	ToAddress               string     `json:"toAddress"`
	Services                serviceDTO `json:"services"`
	BankDetails             string     `json:"bankDetails"`
	AdditionalInfo          string     `json:"additionalInfo"`
	CompanyPromoInfoPhone   string     `json:"companyPromoInfoPhone"`
	CompanyPromoInfoEmail   string     `json:"companyPromoInfoEmail"`
	CompanyPromoInfoWebPage string     `json:"companyPromoInfoWebPage"`
	Tax                     float64    `json:"tax"`
	Discount                float64    `json:"discount"`
	Total                   float64    `json:"total"`
	Subtotal                float64    `json:"subtotal"`
	BalanceDue              float64    `json:"balance_due"`
	PaidDate                string     `json:"paidDate"`
	PaidAmount              float64    `json:"paidAmount"`
	IsExpense               bool       `json:"isExpense"`
}

type serviceDTO struct {
	ID int `json:"id"`
}

type contractDTO struct {
	ID int `json:"id"`
}
