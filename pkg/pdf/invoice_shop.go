package pdf

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v8/pkg/shop/db"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

func shopTaxClassToText(taxClass string) string {
	switch taxClass {
	case db.TaxClassStandard:
		return "19%"
	case db.TaxClassTakeaway:
		return "7%"
	default:
		return "N/A"
	}
}

func InvoiceFromShopOrder(ctx context.Context, order *db.Order, store *db.Store) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Infof("Generate invoice for shop order %s", *order.OrderNumber)

	tax19 := decimal.Zero
	tax7 := decimal.Zero
	items := lo.Map(order.Items, func(item db.OrderItem, _ int) invoiceTemplateItems {
		switch item.TaxClass {
		case db.TaxClassStandard:
			tax19 = tax19.Add(decimal.RequireFromString(item.Tax))
		case db.TaxClassTakeaway:
			tax7 = tax7.Add(decimal.RequireFromString(item.Tax))
		default:
			log.Errorf("Unsupported tax class: %v", item)
		}
		return invoiceTemplateItems{
			Name:     item.Name,
			SKU:      item.SKU,
			Amount:   item.Amount,
			Total:    item.Total,
			Tax:      item.Tax,
			TaxClass: shopTaxClassToText(item.TaxClass),
		}
	})
	props := invoiceTemplateProps{
		StoreName:       store.Name,
		StoreAddress:    store.Address,
		StoreTaxNumber:  store.TaxNumber,
		StoreTelephone:  store.Telephone,
		StoreEmail:      store.Email,
		CustomerName:    fmt.Sprintf("%s %s", order.Checkout.FirstName, order.Checkout.LastName),
		CustomerAddress: fmt.Sprintf("%s %s, %s %s", order.Checkout.AddressLine1, order.Checkout.AddressLine2, order.Checkout.Postcode, order.Checkout.City),
		InvoiceNumber:   *order.InvoiceNumber,
		OrderNumber:     *order.OrderNumber,
		Date:            order.UpdatedAt,
		Items:           items,
		Total:           order.Summary.Total.Value,
		Tax:             order.Summary.Tax.Value,
		Tax19:           tax19.StringFixed(2),
		Tax7:            tax7.StringFixed(2),
	}
	t, err := template.New("invoice").Funcs(
		template.FuncMap{
			"DateTime": dateTimeFromTime,
		},
	).Parse(invoiceTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse invoice template")
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, props); err != nil {
		return nil, errors.Wrap(err, "failed to execute invoice template")
	}
	invoiceHTML, err := io.ReadAll(buf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read invoice template result")
	}
	return NewFromHTML(ctx, string(invoiceHTML))
}
