package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/mongodb"
)

const colCarts = "carts"

func CollectionCarts(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/shop/mongo/database", false)
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colCarts)
}

type Cart struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Items     []CartItem         `bson:"items" json:"items"`
	Secret    string             `bson:"secret" json:"secret"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type CartItem struct {
	ProductID string            `bson:"productId" json:"productId"`
	Options   map[string]string `bson:"optionValues" json:"optionValues"`
	Amount    int               `bson:"amount" json:"amount"`
}
