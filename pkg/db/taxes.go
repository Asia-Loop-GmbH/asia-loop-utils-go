package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/pkg/mongodb"
)

const colTaxes = "taxes"

func CollectionTaxes(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/mongo/database", false)
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colTaxes)
}

type Tax struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	WPID      int                `bson:"id" json:"wpId"`
	Rate      string             `bson:"rate" json:"rate"`
	Name      string             `bson:"name" json:"name"`
	TaxClass  string             `bson:"class" json:"taxClass"` // it's ok to have different names here because we don't provide PATCH request for this entity.
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
