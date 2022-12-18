package servicewoo

import (
	"fmt"
	"time"
)

const (
	metaDataKeyInvoice      = "_wcpdf_invoice_number"
	metaDataKeyDeliveryDate = "jckwds_date_ymd"
	metaDataKeyDeliveryTime = "jckwds_timeslot"
	metaDataKeyStore        = "_AL_STORE"
	feeLineTip              = "Tip (Custom Tip)"
)

type Order struct {
	ID               int             `json:"id"`
	Number           string          `json:"number"`
	DiscountTotal    string          `json:"discount_total"`
	DiscountTax      string          `json:"discount_tax"`
	ShippingTotal    string          `json:"shipping_total"`
	ShippingTax      string          `json:"shipping_tax"`
	CartTax          string          `json:"cart_tax"`
	Total            string          `json:"total"`
	TotalTax         string          `json:"total_tax"`
	PricesIncludeTax bool            `json:"prices_include_tax"`
	CustomerNote     string          `json:"customer_note"`
	Billing          BillingAddress  `json:"billing"`
	Shipping         ShippingAddress `json:"shipping"`
	MetaData         []MetaData      `json:"meta_data"`
	LineItems        []LineItem      `json:"line_items"`
	TaxLines         []TaxLine       `json:"tax_lines"`
	ShippingLines    []ShippingLine  `json:"shipping_lines"`
	FeeLines         []FeeLine       `json:"fee_lines"`
	Refunds          []OrderRefund   `json:"refunds"`
	DateCreatedGMT   string          `json:"date_created_gmt"`
}

func (o *Order) GetInvoiceNumber() string {
	return o.getMetaDataStringWithDefault(metaDataKeyInvoice, "")
}

func (o *Order) GetStoreKey() string {
	return o.getMetaDataStringWithDefault(metaDataKeyStore, "")
}

func (o *Order) GetDeliveryDate() string {
	return o.getMetaDataStringWithDefault(metaDataKeyDeliveryDate, "")
}

func (o *Order) GetDeliveryTime() string {
	return o.getMetaDataStringWithDefault(metaDataKeyDeliveryTime, "")
}

func (o *Order) GetTip() string {
	for _, item := range o.FeeLines {
		if item.Name == feeLineTip {
			return item.Total
		}
	}
	return "0.00"
}

func (o *Order) getMetaDataStringWithDefault(key, defaultValue string) string {
	v, err := o.getMetaDataString(key)
	if err != nil {
		return defaultValue
	}
	return v
}

func (o *Order) getMetaDataString(key string) (string, error) {
	for _, item := range o.MetaData {
		if item.Key == key {
			switch v := item.Value.(type) {
			case string:
				return v, nil
			default:
				return "", fmt.Errorf("expect string value, got: %v", item)
			}
		}
	}
	return "", fmt.Errorf("key [%s] not found in: %v", key, o)
}

func (o *Order) GetShippingMethod() string {
	if len(o.ShippingLines) > 0 {
		return o.ShippingLines[0].MethodID
	} else {
		return "no_shipping"
	}
}

func (o *Order) GetDateCreated() (time.Time, error) {
	loc, err := time.LoadLocation("GMT")
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation("2006-01-02T15:04:05", o.DateCreatedGMT, loc)
}
