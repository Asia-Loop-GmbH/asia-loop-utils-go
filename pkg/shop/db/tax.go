package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v2/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v2/pkg/mongodb"
)

const colTax = "taxes"

func CollectionTaxes(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/shop/mongo/database", false)
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
}
