package api

type SearchCustomerRequest struct {
	Text         *string `json:"text"`
	FirstName    *string `json:"firstName"`
	LastName     *string `json:"lastName"`
	AddressLine1 *string `json:"addressLine1"`
	Postcode     *string `json:"postcode"`
	Telephone    *string `json:"telephone"`
	Email        *string `json:"email"`
	Limit        *int64  `json:"limit"`
}
