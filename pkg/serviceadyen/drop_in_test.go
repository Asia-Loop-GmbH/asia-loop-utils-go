package serviceadyen_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v3/pkg/serviceadyen"
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v3/pkg/shop/cart"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
)

func TestNewDropInPayment_Success(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	response, err := serviceadyen.NewDropInPayment(
		ctx,
		&cart.PublicCart{
			ID: primitive.NewObjectID(),
			Summary: cart.PublicCartSummary{
				Total: cart.TotalSummary{Value: "12.34"},
			},
		},
		"http://localhost:3000",
	)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}
