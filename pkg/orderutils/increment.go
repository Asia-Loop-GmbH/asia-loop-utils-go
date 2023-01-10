package orderutils

import (
	"context"
	"fmt"
	"time"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v2/pkg/db"
)

const (
	incrementOrderInvoiceKey           = "ORDER_INVOICE"
	incrementOrderInvoiceKeyLieferando = "ORDER_INVOICE_LIEFERANDO"
)

func NextOrderInvoice(ctx context.Context) (*string, error) {
	next, err := db.Next(ctx, incrementOrderInvoiceKey)
	if err != nil {
		return nil, err
	}
	prefix, err := orderInvoicePrefix()
	if err != nil {
		return nil, err
	}
	full := fmt.Sprintf("P%s-%07d", *prefix, next)
	return &full, nil
}

func NextOrderInvoiceLieferando(ctx context.Context) (*string, error) {
	next, err := db.Next(ctx, incrementOrderInvoiceKeyLieferando)
	if err != nil {
		return nil, err
	}
	prefix, err := orderInvoicePrefix()
	if err != nil {
		return nil, err
	}
	full := fmt.Sprintf("L%s-%07d", *prefix, next)
	return &full, nil
}

func orderInvoicePrefix() (*string, error) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return nil, err
	}
	prefix := time.Now().In(loc).Format("200601")
	return &prefix, nil
}
