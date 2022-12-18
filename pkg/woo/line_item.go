package woo

type LineItem struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	ProductID   int               `json:"product_id"`
	VariationID *int              `json:"variation_id"`
	Quantity    int               `json:"quantity"`
	TaxClass    string            `json:"tax_class"`
	SubTotal    string            `json:"subtotal"`
	SubTotalTax string            `json:"subtotal_tax"`
	Total       string            `json:"total"`
	TotalTax    string            `json:"total_tax"`
	Taxes       []Tax             `json:"taxes"`
	MetaData    []MetaData        `json:"meta_data"`
	SKU         string            `json:"sku"`
	Price       float32           `json:"price"`
	Categories  []ProductCategory `json:"_al_categories"`
}
