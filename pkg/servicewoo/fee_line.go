package servicewoo

type FeeLine struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	TaxClass  string     `json:"tax_class"`
	TaxStatus string     `json:"tax_status"`
	Amount    string     `json:"amount"`
	Total     string     `json:"total"`
	TotalTax  string     `json:"total_tax"`
	MetaData  []MetaData `json:"meta_data"`
}
