package servicewoo

type ProductVariationAttribute struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Option string `json:"option"`
}

type ProductVariation struct {
	ID           int                         `json:"id"`
	Price        string                      `json:"price"`
	RegularPrice string                      `json:"regular_price"`
	SalePrice    string                      `json:"sale_price"`
	Attributes   []ProductVariationAttribute `json:"attributes"`
}
