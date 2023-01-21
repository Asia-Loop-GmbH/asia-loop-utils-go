package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/mongodb"
)

type Order struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	StoreKey      string             `bson:"storeKey" json:"storeKey"`
	IsPickup      bool               `bson:"isPickup" json:"isPickup"`
	Items         []OrderItem        `bson:"items" json:"items"`
	Summary       OrderSummary       `bson:"summary" json:"summary"`
	Secret        string             `bson:"secret" json:"secret"`
	Paid          bool               `bson:"paid" json:"paid"`
	InvoiceNumber *string            `bson:"invoiceNumber,omitempty" json:"invoiceNumber,omitempty"`
	OrderNumber   *string            `bson:"orderNumber,omitempty" json:"orderNumber,omitempty"`
	Payment       *Payment           `bson:"payment" json:"payment,omitempty"`
	Checkout      *CartCheckout      `bson:"checkout,omitempty" json:"checkout,omitempty"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type OrderSummary struct {
	Total  TotalSummary `bson:"total" json:"total"`
	Tax    TotalSummary `bson:"tax" json:"tax"`
	Net    TotalSummary `bson:"net" json:"net"`
	Saving string       `bson:"saving" json:"saving"`
}

type TotalSummary struct {
	Value  string            `bson:"value" json:"value"`
	Values map[string]string `bson:"values" json:"values"` // Values are values grouped by tax classes
}

type OrderItem struct {
	CartItem
	SKU        string   `bson:"sku" json:"sku"`
	Name       string   `bson:"name" json:"name"`
	Categories []string `bson:"categories" json:"categories"`
	UnitPrice  string   `bson:"unitPrice" json:"unitPrice"`
	Total      string   `bson:"total" json:"total"`
	Tax        string   `bson:"tax" json:"tax"`
	Net        string   `bson:"net" json:"net"`
	Saving     string   `bson:"saving" json:"saving"`
	TaxClass   string   `bson:"taxClass" json:"taxClass"`
}

const colOrders = "orders"

func CollectionOrders(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/shop/mongo/database", false)
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colOrders)
}
