package orderutils_test

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v3/pkg/orderutils"
	commoncontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
)

func TestNextOrderInvoice(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.Background(), commoncontext.FieldStage, "dev")
	next, err := orderutils.NextOrderInvoice(ctx)
	assert.NoError(t, err)
	log.Printf("%s", *next)
}

func TestNextOrderInvoiceLieferando(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.Background(), commoncontext.FieldStage, "dev")
	next, err := orderutils.NextOrderInvoiceLieferando(ctx)
	assert.NoError(t, err)
	log.Printf("%s", *next)
}

func TestNextShopOrderInvoice(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.Background(), commoncontext.FieldStage, "dev")
	next, err := orderutils.NextShopOrderInvoice(ctx)
	assert.NoError(t, err)
	log.Printf("%s", *next)
}
