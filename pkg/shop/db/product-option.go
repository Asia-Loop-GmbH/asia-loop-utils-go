package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v2/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v2/pkg/mongodb"
)

const colProductOptions = "product-options"

func CollectionProductOptions(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/shop/mongo/database", false)
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colProductOptions)
}

type ProductOption struct {
	ID          primitive.ObjectID   `bson:"_id" json:"id"`
	Name        string               `bson:"name" json:"name"`
	DisplayName string               `bson:"displayName" json:"displayName"`
	Values      []ProductOptionValue `bson:"values" json:"values"`
	CreatedAt   time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time            `bson:"updatedAt" json:"updatedAt"`
}

type ProductOptionValue struct {
	Name         string   `bson:"name" json:"name"`
	DisplayName  string   `bson:"displayName" json:"displayName"`
	DisabledIn   []string `bson:"disabledIn" json:"disabledIn"`
	OutOfStockIn []string `bson:"outOfStockIn" json:"outOfStockIn"`
}
