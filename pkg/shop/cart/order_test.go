package cart

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v9/pkg/shop/db"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
)

func TestCreateOrder(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	_ = os.Setenv("APP", "admin")
	ctx := context.WithValue(context.Background(), mycontext.FieldStage, "prod")

	colCarts, err := db.CollectionCarts(ctx)
	assert.NoError(t, err)
	cartID, err := primitive.ObjectIDFromHex("669a956a4bdad7e1c9e7ad1d")
	assert.NoError(t, err)
	find := colCarts.FindOne(ctx, bson.M{"_id": cartID})
	shoppingCart := new(db.Cart)
	err = find.Decode(shoppingCart)
	assert.NoError(t, err)
	order, err := ToOrder(ctx, shoppingCart)
	assert.NoError(t, err)

	_, err = CreateOrder(ctx, shoppingCart, order.Summary.Total.Value)
	assert.NoError(t, err)
}
