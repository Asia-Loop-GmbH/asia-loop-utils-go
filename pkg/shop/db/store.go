package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/mongodb"
)

const colStores = "stores"

func CollectionStores(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/shop/mongo/database", false)
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colStores)
}

type Store struct {
	ID                   primitive.ObjectID `bson:"_id" json:"id"`
	Email                string             `bson:"email" json:"email"`
	Telephone            string             `bson:"telephone" json:"telephone"`
	Name                 string             `bson:"name" json:"name"`
	Key                  string             `bson:"key" json:"key"`
	Address              string             `bson:"address" json:"address"`
	Owner                string             `bson:"owner" json:"owner"`
	BusinessRegistration string             `bson:"businessRegistration" json:"businessRegistration"`
	TaxNumber            string             `bson:"taxNumber" json:"taxNumber"`
	MBW                  map[string]string  `bson:"mbw" json:"mbw"`
	CreatedAt            time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt            time.Time          `bson:"updatedAt" json:"updatedAt"`
}
