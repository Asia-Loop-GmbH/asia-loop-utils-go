package cart_test

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v2/pkg/shop/cart"
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v2/pkg/shop/db"
)

var (
	taxes = []db.Tax{
		{
			ID:          primitive.NewObjectID(),
			Name:        "takeaway",
			DisplayName: "MwSt. 7%",
			Rate:        "0.07",
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		},
	}
)

func TestDecimal(t *testing.T) {
	assert.Equal(t, "1.22", decimal.NewFromFloat(1.224).Round(2).StringFixed(2))
	assert.Equal(t, "1.23", decimal.NewFromFloat(1.227).Round(2).StringFixed(2))
	assert.Equal(t, "1.23", decimal.NewFromFloat(1.225).Round(2).StringFixed(2))
}

func TestCalculate(t *testing.T) {
	id := primitive.NewObjectID()
	secret := "some-secret"
	updated := time.Now()
	created := time.Now()
	sku := "some-sku"
	name := "product name"
	price := "12.34"
	amount := 4
	expectedTotal := "49.36"
	expectedTax := "3.24"
	expectedNet := "46.12"
	expectedTaxClass := "takeaway"
	cartItem := db.CartItem{
		ProductID: id.String(),
		Options:   nil,
		Amount:    amount,
	}

	categories := []string{"category1", "category2"}
	products := []db.Product{
		{
			ID:   id,
			SKU:  sku,
			Name: name,
			Price: db.CustomizablePrice{
				Value:        price,
				CustomValues: nil,
			},
			TaxClassStandard: "standard",
			TaxClassTakeAway: "takeaway",
			Categories:       categories,
			Options:          nil,
			Variations:       nil,
		},
	}
	shoppingCart := db.Cart{
		ID:        id,
		IsPickup:  false,
		Secret:    secret,
		CreatedAt: created,
		UpdatedAt: updated,
		Items:     []db.CartItem{cartItem},
	}
	publicCart, err := cart.Calculate(context.TODO(), &shoppingCart, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, id, publicCart.ID)
	assert.Equal(t, secret, publicCart.Secret)
	assert.Equal(t, updated, publicCart.UpdatedAt)
	assert.Equal(t, created, publicCart.CreatedAt)
	assert.Equal(t, false, publicCart.IsPickup)

	assert.Equal(t, cart.PublicCartItem{
		CartItem:   cartItem,
		SKU:        sku,
		Categories: categories,
		UnitPrice:  price,
		Total:      expectedTotal,
		Tax:        expectedTax,
		Net:        expectedNet,
		TaxClass:   expectedTaxClass,
	}, publicCart.Items[0])

	assert.Equal(t, cart.PublicCartSummary{
		Total:    "49.36",
		TotalTax: "3.24",
		TotalNet: "46.12",
		Taxes: map[string]string{
			"takeaway": "3.24",
		},
	}, publicCart.Summary)
}
