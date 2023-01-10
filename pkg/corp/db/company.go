package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/mongodb"
)

const colCompany = "companies"

func CollectionCompanies(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/corp/mongo/database", false)
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colCompany)
}

type Company struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	Key             string             `bson:"key" json:"key"`
	Name            string             `bson:"name" json:"name"`
	Address1        string             `bson:"address1" json:"address1"`
	Address2        string             `bson:"address2" json:"address2"`
	Postcode        string             `bson:"postcode" json:"postcode"`
	City            string             `bson:"city" json:"city"`
	ContactName     string             `bson:"contactName" json:"contactName"`
	ContactEmail    string             `bson:"contactEmail" json:"contactEmail"`
	ContactPhone    string             `bson:"contactPhone" json:"contactPhone"`
	PartnerStoreKey string             `bson:"partnerStoreKey" json:"partnerStoreKey"`
}
