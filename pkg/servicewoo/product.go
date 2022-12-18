package servicewoo

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Image struct {
	Source string `json:"src"`
}

type ProductAttribute struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Position  int      `json:"position"`
	Variation bool     `json:"variation"`
	Visible   bool     `json:"visible"`
	Options   []string `json:"options"`
}

type Product struct {
	ID           int                `json:"id"`
	Name         string             `json:"name"`
	Type         string             `json:"type"`
	Status       string             `json:"status"`
	SKU          string             `json:"sku"`
	Price        string             `json:"price"`
	RegularPrice string             `json:"regular_price"`
	SalePrice    string             `json:"sale_price"`
	TaxClass     string             `json:"tax_class"`
	Categories   []Category         `json:"categories"`
	Images       []Image            `json:"images"`
	Attributes   []ProductAttribute `json:"attributes"`
	Variations   []int              `json:"variations"`
}
