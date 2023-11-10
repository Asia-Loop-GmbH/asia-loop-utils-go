package deliverect

type WebhookProductsResponse struct {
	AccountID   string         `json:"accountId"`
	LocationID  string         `json:"locationId"`
	Products    []Product      `json:"products"`
	Categories  *[]Category    `json:"categories,omitempty"`
	PriceLevels *[]PriceLevels `json:"priceLevels,omitempty"`
}
