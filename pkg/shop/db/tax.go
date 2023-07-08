package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mongodb"
)

const (
	TaxClassStandard = "standard"
	TaxClassTakeaway = "takeaway"
	TaxClassZero     = "kein"
)

const colTax = "taxes"

func CollectionTaxes(ctx context.Context) (*mongo.Collection, error) {
	database, err := secretsmanager.GetParameter(ctx, "/shop/mongo/database")
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colTax)
}

type Tax struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	DisplayName string             `bson:"displayName" json:"displayName"`
	Rate        string             `bson:"rate" json:"rate"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
