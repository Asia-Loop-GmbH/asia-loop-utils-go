package woo

type TaxLine struct {
	ID               int        `json:"id"`
	RateCode         string     `json:"rate_code"`
	RateID           int        `json:"rate_id"`
	Label            string     `json:"label"`
	Compound         bool       `json:"compound"`
	TaxTotal         string     `json:"tax_total"`
	ShippingTaxTotal string     `json:"shipping_tax_total"`
	RatePercent      int        `json:"rate_percent"`
	MetaData         []MetaData `json:"meta_data"`
}
