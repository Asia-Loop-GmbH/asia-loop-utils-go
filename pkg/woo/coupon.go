package woo

// Coupon : "discount_type", "date_expires", "individual_use",... | currently assume that only one type is used and no expires
// TODO: add logic to differentiate types
type Coupon struct {
	ID     int    `json:"id,omitempty"`
	Code   string `json:"code,omitempty"`
	Amount string `json:"amount"`
}
