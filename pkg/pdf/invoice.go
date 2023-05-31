package pdf

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v8/pkg/db"
	shopdb "github.com/asia-loop-gmbh/asia-loop-utils-go/v8/pkg/shop/db"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
)

func InvoiceFromShopOrder(ctx context.Context, order *shopdb.Order, store *shopdb.Store) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Infof("Generate invoice for shop order %s", *order.OrderNumber)

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
		Items: lo.Map(order.Items, func(item shopdb.OrderItem, _ int) invoiceTemplateItems {
			return invoiceTemplateItems{
				Name:   item.Name,
				SKU:    item.SKU,
				Amount: item.Amount,
				Total:  item.Total,
			}
		}),
		Total: order.Summary.Total.Value,
		Tax:   order.Summary.Tax.Value,
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

func InvoiceFromAdminOrder(ctx context.Context, order *db.Order, store *db.Store, customer *db.Customer) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Infof("Generate invoice for admin order %s", order.OrderNumber)

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
		Items: lo.Map(order.Items, func(item db.OrderItem, _ int) invoiceTemplateItems {
			return invoiceTemplateItems{
				Name:   item.Name,
				SKU:    item.SKU,
				Amount: item.Quantity,
				Total:  item.Total,
			}
		}),
		Total: order.Total,
		Tax:   order.Tax,
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

func dateTimeFromTime(x time.Time) string {
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return "N/A"
	}
	return x.In(location).Format("02.01.2006 15:04:05")
}
