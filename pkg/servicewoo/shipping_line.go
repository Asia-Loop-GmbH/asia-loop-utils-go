package servicewoo

type ShippingLine struct {
	ID          int        `json:"id"`
	MethodTitle string     `json:"method_title"`
	MethodID    string     `json:"method_id"`
	Total       string     `json:"total"`
	TotalTax    string     `json:"total_tax"`
	Taxes       []Tax      `json:"taxes"`
	MetaData    []MetaData `json:"meta_data"`
}
