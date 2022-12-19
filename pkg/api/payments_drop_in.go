package api

type PaymentDropInRequest struct {
	ReturnURL string `json:"returnUrl"`
}

type PaymentDropInResponse struct {
	ID          string `json:"id"`
	SessionData string `json:"sessionData"`
}
