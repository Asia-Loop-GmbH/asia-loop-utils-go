package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mongodb"
)

const colGlobalConfig = "globalconfigs"

func CollectionGlobalConfig(ctx context.Context) (*mongo.Collection, error) {
	database, err := secretsmanager.GetParameter(ctx, "/mongo/database")
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colGlobalConfig)
}

type GlobalConfig struct {
	ID                           primitive.ObjectID `bson:"_id" json:"id"`
	ProductAttributeOutOfStockIn int                `bson:"productAttributeOutOfStockInId" json:"productAttributeOutOfStockInId"`
	ProductAttributePfandId      int                `bson:"productAttributePfandId" json:"productAttributePfandId"`
	PusherAPIKey                 string             `json:"pusherApiKey"`
	DeliverectWebhookSecret      string             `json:"deliverectWebhookSecret"`
}
