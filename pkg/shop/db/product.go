package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/mongodb"
)

const colProducts = "products"

func CollectionProducts(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/shop/mongo/database", false)
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colProducts)
}

type Product struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	SKU              string             `bson:"sku" json:"sku"`
	Name             string             `bson:"name" json:"name"`
	Price            CustomizablePrice  `bson:"price" json:"price"`
	TaxClassStandard string             `bson:"taxClassStandard" json:"taxClassStandard"`
	TaxClassTakeAway string             `bson:"taxClassTakeAway" json:"taxClassTakeAway"`
	Categories       []string           `bson:"categories" json:"categories"`
	Images           []Image            `bson:"images" json:"images"`
	Options          []ProductOption    `bson:"options" json:"options"`
	Variations       []ProductVariation `bson:"variations" json:"variations"`
	DisabledIn       []string           `bson:"disabledIn" json:"disabledIn"`
	OutOfStockIn     []string           `bson:"outOfStockIn" json:"outOfStockIn"`
	Description      string             `bson:"description" json:"description"`
	Allergens        []string           `bson:"allergens" json:"allergens"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type ProductVariation struct {
	SKU     string                   `bson:"sku" json:"sku"`
	Price   CustomizablePrice        `bson:"price" json:"price"`
	Options []ProductVariationOption `bson:"options" json:"options"`
}

type ProductVariationOption struct {
	Name  string `bson:"name" json:"name"`
	Value string `bson:"value" json:"value"`
}
