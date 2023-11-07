package db

import (
	"context"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mongodb"
)

const colProducts = "products"

func CollectionProducts(ctx context.Context) (*mongo.Collection, error) {
	database, err := secretsmanager.GetParameter(ctx, "/shop/mongo/database")
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
	Options          []string           `bson:"options" json:"options"`
	Variations       []ProductVariation `bson:"variations" json:"variations"`
	DisabledIn       []string           `bson:"disabledIn" json:"disabledIn"`
	OutOfStockIn     []string           `bson:"outOfStockIn" json:"outOfStockIn"`
	Description      string             `bson:"description" json:"description"`
	Allergens        []string           `bson:"allergens" json:"allergens"`
	IsGiftCard       bool               `bson:"isGiftCard" json:"isGiftCard"`
	IsStar           bool               `bson:"isStar" json:"isStar"`
	IsNew            bool               `bson:"isNew" json:"isNew"`
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

func (p *Product) GetPrice(storeKey string, selectedOptions map[string]string) string {
	if p.IsGiftCard {
		return p.Price.Value
	}

	variation, ok := lo.Find(p.Variations, func(variation ProductVariation) bool {
		return lo.EveryBy(variation.Options, func(opt ProductVariationOption) bool {
			return opt.Value == selectedOptions[opt.Name]
		})
	})
	if ok {
		return variation.Price.GetPrice(storeKey)
	}
	return p.Price.GetPrice(storeKey)
}

func (p *Product) IsAvailableInStore(storeKey string, selectedOptions map[string]string, options []ProductOption) bool {
	if p.IsGiftCard {
		return false
	}
	if lo.Contains(p.DisabledIn, storeKey) || lo.Contains(p.OutOfStockIn, storeKey) {
		return false
	}

	for selectedName, selectedValue := range selectedOptions {
		for _, option := range options {
			if selectedName == option.Name {
				for _, value := range option.Values {
					if value.Name == selectedValue {
						if lo.Contains(value.DisabledIn, storeKey) || lo.Contains(value.OutOfStockIn, storeKey) {
							return false
						}
					}
				}
			}
		}
	}

	return true
}
