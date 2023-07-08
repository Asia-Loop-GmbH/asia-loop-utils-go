package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mongodb"
)

const colProductOptions = "product-options"

func CollectionProductOptions(ctx context.Context) (*mongo.Collection, error) {
	database, err := secretsmanager.GetParameter(ctx, "/shop/mongo/database")
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
	PrintName    string   `bson:"printName" json:"printName"`
	Image        *Image   `bson:"image" json:"image"`
	DisabledIn   []string `bson:"disabledIn" json:"disabledIn"`
	OutOfStockIn []string `bson:"outOfStockIn" json:"outOfStockIn"`
	Allergens    []string `bson:"allergens" json:"allergens"`
}
