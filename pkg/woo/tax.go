package woo

type Tax struct {
	ID       int    `json:"id"`
	Total    string `json:"total"`
	SubTotal string `json:"subtotal"`
}
