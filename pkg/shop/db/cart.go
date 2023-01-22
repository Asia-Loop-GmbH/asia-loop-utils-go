package db

import (
	"context"
	"time"

	"github.com/adyen/adyen-go-api-library/v6/src/checkout"
	"github.com/adyen/adyen-go-api-library/v6/src/notification"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/mongodb"
)

const colCarts = "carts"

func CollectionCarts(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/shop/mongo/database", false)
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colCarts)
}

type Cart struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	User      *string            `bson:"user,omitempty" json:"user,omitempty"`
	StoreKey  string             `bson:"storeKey" json:"storeKey"`
	IsPickup  bool               `bson:"isPickup" json:"isPickup"`
	Items     []CartItem         `bson:"items" json:"items"`
	Secret    string             `bson:"secret" json:"secret"`
	Payments  []Payment          `bson:"payments" json:"payments"`
	Paid      bool               `bson:"paid" json:"paid"`
	Checkout  *CartCheckout      `bson:"checkout,omitempty" json:"checkout,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type CartCheckout struct {
	FirstName    string     `bson:"firstName" json:"firstName"`
	LastName     string     `bson:"lastName" json:"lastName"`
	AddressLine1 string     `bson:"addressLine1" json:"addressLine1"`
	AddressLine2 string     `bson:"addressLine2" json:"addressLine2"`
	City         string     `bson:"city" json:"city"`
	Postcode     string     `bson:"postcode" json:"postcode"`
	Telephone    string     `bson:"telephone" json:"telephone"`
	Email        string     `bson:"email" json:"email"`
	Note         string     `bson:"note" json:"note"`
	Date         string     `bson:"date" json:"date"`
	Slot         string     `bson:"slot" json:"slot"`
	Begin        *time.Time `bson:"begin,omitempty" json:"begin,omitempty"`
}

type Payment struct {
	Session     checkout.CreateCheckoutSessionResponse `bson:"session" json:"session"`
	Environment string                                 `bson:"environment" json:"environment"`
	Client      string                                 `bson:"client" json:"client"`
	Events      []notification.NotificationRequestItem `bson:"events" json:"events"`
}

type CartItem struct {
	ProductID string            `bson:"productId" json:"productId"`
	Options   map[string]string `bson:"optionValues" json:"optionValues"`
	Amount    int               `bson:"amount" json:"amount"`
}
