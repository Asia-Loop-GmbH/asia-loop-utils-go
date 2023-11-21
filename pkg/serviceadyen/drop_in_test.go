package serviceadyen_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v9/pkg/serviceadyen"
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v9/pkg/shop/db"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
)

func TestNewDropInPayment_Success(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	response, err := serviceadyen.NewDropInPayment(
		ctx,
		&db.Order{
			ID: primitive.NewObjectID(),
			Summary: db.OrderSummary{
				Total: db.TotalSummary{Value: "12.34"},
			},
			Items: []db.OrderItem{
				{
					SKU:      "SKU",
					Name:     "Name",
					Total:    "12.34",
					Tax:      "2.34",
					Net:      "10.00",
					TaxClass: "700",
				},
			},
		},
		&db.CartCheckout{
			Email: "lenamtruong@gmail.com",
		},
		"http://localhost:3000",
	)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}
