Return ONE JSON OBJECT.
- All fields must be mandatorily present in the response.
- The service list will be nested in the contract always.
- For the currency symbol surrounded by words such as "total" "sub total", etc.
- The CurrencySymbol Must be a 3 letter code of the likes of "MDL", "USD", "EUR", etc.
- Frequency MUST be classified as  "monthly" by default, otherwise deduced from sections such as "Term and Termination".
- Amount is the quantity in termsm of money that describes one service payment.
- Taxes could be identified by "Fees" or "VTA"
- Quantity is the amount of units.
- PriceUnit is the Amount paid per ONE Unit.
- Item is the title of the Service
- Time Must be in Golang time.Time compatible format
- Use the given legal document as context.
```go
type serviceDTO struct {
	Item                    string     `json:"item"`
	Quantity                int     `json:"quantity"`
	Frequency 			string	       `json:"frequency"`
	PriceUnit				float64 	   `json:"priceUnit"`
	Taxes					float64	   `json:"taxes"`
	Amount 					float64	   `json:"amount"`
}

type contractDTO struct {
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
	Taxes					float64     `json:"taxes"`
}
```
