Return ONE JSON OBJECT, all fields mandatory.
Use the given legal document as context.
For the currency symbol use the one at total.
Frequency will be either "yearly", "monthly", "weekly" and deduced from the terms.

```go
type serviceDTO struct {
	ID 						int 	   `json:"id"`
	Item                    string     `json:"item"`
	Quantity                int     `json:"quantity"`
	PriceUnit				float64 	   `json:"priceUnit"`
	Taxes					float64	   `json:"taxes"`
	Amount 					float64	   `json:"amount"`
}

type contractDTO struct {
	ID                      int        `json:"id"`
	UserID                  int        `json:"userId"`
	URLHash                 string     `json:"urlHash"`
	DocID                   string     `json:"docId"`
	CurrencySymbol          string     `json:"currencySymbol"`
	LanguageCode            string     `json:"languageCode"`
	IssueDate               time.Time  `json:"issueDate"`
	DueDate                 time.Time  `json:"dueDate"`
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
	Frequency 			string	       `json:"frequency"`
	Subtotal 				float64     `json:"amount"`
	Taxes					float64     `json:"taxes"`
	Total  					float64     `json:"total"`
}
```
