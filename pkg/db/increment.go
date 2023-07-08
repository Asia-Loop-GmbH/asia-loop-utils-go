package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mongodb"
)

const colIncrement = "increments"

func collectionIncrement(ctx context.Context) (*mongo.Collection, error) {
	database, err := secretsmanager.GetParameter(ctx, "/mongo/database")
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colIncrement)
}

type increment struct {
	ID        primitive.ObjectID `bson:"_id"`
	Key       string             `bson:"key"`
	Value     int64              `bson:"value"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
