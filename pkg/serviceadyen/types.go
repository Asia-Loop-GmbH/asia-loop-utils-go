package serviceadyen

const (
	protocolVersion        = "3.0"
	messageTypeRequest     = "Request"
	messageClassService    = "Service"
	messageCategoryPayment = "Payment"
	currencyEUR            = "EUR"
	countryDE              = "DE"
	accountECOM            = "AsiaLoopECOM"
	livePrefix             = "9567a7305fc929ca-AsiaLoopGmbH"
)

type MessageHeader struct {
	ProtocolVersion string `json:"ProtocolVersion"`
	MessageClass    string `json:"MessageClass"`
	MessageCategory string `json:"MessageCategory"`
	MessageType     string `json:"MessageType"`
	SaleID          string `json:"SaleID"`
	ServiceID       string `json:"ServiceID"`
	POIID           string `json:"POIID"`
}

type SaleTransactionID struct {
	TransactionID string `json:"TransactionID"`
	TimeStamp     string `json:"TimeStamp"`
}

type SaleData struct {
	SaleTransactionID SaleTransactionID `json:"SaleTransactionID"`
}

type AmountsReq struct {
	Currency        string  `json:"Currency"`
	RequestedAmount float32 `json:"RequestedAmount"`
}

type PaymentTransaction struct {
	AmountsReq AmountsReq `json:"AmountsReq"`
}

type PaymentRequest struct {
	SaleData           SaleData           `json:"SaleData"`
	PaymentTransaction PaymentTransaction `json:"PaymentTransaction"`
}

type SaleToPOIRequest struct {
	MessageHeader  MessageHeader  `json:"MessageHeader"`
	PaymentRequest PaymentRequest `json:"PaymentRequest"`
}

type TerminalAPIRequest struct {
	SaleToPOIRequest SaleToPOIRequest `json:"SaleToPOIRequest"`
}
