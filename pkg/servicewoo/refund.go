package servicewoo

import "time"

type Refund struct {
	ID             int        `json:"id"`
	DateCreated    string     `json:"date_created"`
	DateCreatedGMT string     `json:"date_created_gmt"`
	Amount         string     `json:"amount"`
	Reason         string     `json:"reason"`
	MetaData       []MetaData `json:"meta_data"`
	LineItems      []LineItem `json:"line_items"`
	TaxLines       []TaxLine  `json:"tax_lines"`
	FeeLines       []FeeLine  `json:"fee_lines"`
}

type OrderRefund struct {
	ID     int    `json:"id"`
	Reason string `json:"reason"`
	Total  string `json:"total"`
}

func (re *Refund) GetDateCreated() (time.Time, error) {
	loc, err := time.LoadLocation("GMT")
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation("2006-01-02T15:04:05", re.DateCreatedGMT, loc)
}
