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

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v8/pkg/db"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

func adminTaxClassToText(taxClass db.TaxClass) string {
	switch taxClass {
	case db.TaxClassStandard:
		return "19%"
	case db.TaxClassReduced:
		return "7%"
	default:
		return "N/A"
	}
}

func InvoiceFromAdminOrder(ctx context.Context, order *db.Order, store *db.Store, customer *db.Customer) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Infof("Generate invoice for admin order %s", order.OrderNumber)

	tax19 := decimal.Zero
	tax7 := decimal.Zero
	items := lo.Map(order.Items, func(item db.OrderItem, _ int) invoiceTemplateItems {
		if item.TaxClass == db.TaxClassReduced {
			tax7 = tax7.Add(decimal.RequireFromString(item.Tax))
		}

		if item.TaxClass == db.TaxClassStandard {
			tax19 = tax19.Add(decimal.RequireFromString(item.Tax))
		}

		return invoiceTemplateItems{
			Name:     item.Name,
			SKU:      item.SKU,
			Amount:   item.Quantity,
			Total:    item.Total,
			TaxClass: adminTaxClassToText(item.TaxClass),
			Tax:      item.Tax,
		}
	})
	props := invoiceTemplateProps{
		StoreName:       store.Name,
		StoreAddress:    store.Address,
		StoreTaxNumber:  store.Tax,
		StoreTelephone:  store.Telephone,
		StoreEmail:      store.Email,
		CustomerName:    fmt.Sprintf("%s %s", customer.FirstName, customer.LastName),
		CustomerAddress: fmt.Sprintf("%s %s, %s %s", customer.AddressLine1, customer.AddressLine2, customer.Postcode, customer.City),
		InvoiceNumber:   order.InvoiceNumber,
		OrderNumber:     order.OrderNumber,
		Date:            order.UpdatedAt,
		Items:           items,
		Total:           order.Total,
		Tax:             order.Tax,
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
