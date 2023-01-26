package serviceadyen_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v7/pkg/serviceadyen"
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v7/pkg/shop/db"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
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
		},
		"http://localhost:3000",
	)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}
