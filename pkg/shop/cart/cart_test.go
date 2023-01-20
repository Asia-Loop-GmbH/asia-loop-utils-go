package cart_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/adyen/adyen-go-api-library/v6/src/checkout"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v4/pkg/shop/cart"
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v4/pkg/shop/db"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
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

func TestCalculate_IgnoreExpired(t *testing.T) {
	id := primitive.NewObjectID()
	sku := "some-sku"
	name := "product name"
	price := "12.34"
	expectedUnitPrice := "9.87"
	amount := 4
	expectedTotal := "39.48"
	expectedTax := "2.60"
	expectedNet := "36.88"
	expectedTaxClass := "takeaway"
	cartItem := db.CartItem{
		ProductID: id.Hex(),
		Options:   nil,
		Amount:    amount,
	}
	expectedItemSaving := "9.88"
	expectedCartSaving := "9.88"

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
	p1 := db.Payment{
		Session: checkout.CreateCheckoutSessionResponse{
			ExpiresAt: time.Now().Add(-time.Hour),
			Amount: checkout.Amount{
				Currency: "EUR",
				Value:    3948,
			},
		},
		Environment: "env1",
		Client:      "client1",
	}
	shoppingCart := db.Cart{
		ID:       id,
		IsPickup: true,
		Items:    []db.CartItem{cartItem},
		Payments: []db.Payment{p1},
	}
	publicCart, err := cart.Calculate(context.TODO(), &shoppingCart, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, true, publicCart.IsPickup)
	assert.Nil(t, publicCart.Payment)

	assert.Equal(t, db.OrderItem{
		CartItem:   cartItem,
		SKU:        sku,
		Name:       name,
		Categories: categories,
		UnitPrice:  expectedUnitPrice,
		Total:      expectedTotal,
		Tax:        expectedTax,
		Net:        expectedNet,
		Saving:     expectedItemSaving,
		TaxClass:   expectedTaxClass,
	}, publicCart.Items[0])

	assert.Equal(t, db.OrderSummary{
		Total: db.TotalSummary{
			Value:  expectedTotal,
			Values: map[string]string{db.TaxClassTakeaway: expectedTotal},
		},
		Tax: db.TotalSummary{
			Value:  expectedTax,
			Values: map[string]string{db.TaxClassTakeaway: expectedTax},
		},
		Net: db.TotalSummary{
			Value:  expectedNet,
			Values: map[string]string{db.TaxClassTakeaway: expectedNet},
		},
		Saving: expectedCartSaving,
	}, publicCart.Summary)
}

func TestCalculate_IgnoreTotalNotMatch(t *testing.T) {
	id := primitive.NewObjectID()
	sku := "some-sku"
	name := "product name"
	price := "12.34"
	expectedUnitPrice := "9.87"
	amount := 4
	expectedTotal := "39.48"
	expectedTax := "2.60"
	expectedNet := "36.88"
	expectedTaxClass := "takeaway"
	cartItem := db.CartItem{
		ProductID: id.Hex(),
		Options:   nil,
		Amount:    amount,
	}
	expectedItemSaving := "9.88"
	expectedCartSaving := "9.88"

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
	p1 := db.Payment{
		Session: checkout.CreateCheckoutSessionResponse{
			ExpiresAt: time.Now().Add(time.Hour),
			Amount: checkout.Amount{
				Currency: "EUR",
				Value:    1234,
			},
		},
		Environment: "env1",
		Client:      "client1",
	}
	shoppingCart := db.Cart{
		ID:       id,
		IsPickup: true,
		Items:    []db.CartItem{cartItem},
		Payments: []db.Payment{p1},
	}
	publicCart, err := cart.Calculate(context.TODO(), &shoppingCart, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, true, publicCart.IsPickup)
	assert.Nil(t, publicCart.Payment)

	assert.Equal(t, db.OrderItem{
		CartItem:   cartItem,
		SKU:        sku,
		Name:       name,
		Categories: categories,
		UnitPrice:  expectedUnitPrice,
		Total:      expectedTotal,
		Tax:        expectedTax,
		Net:        expectedNet,
		Saving:     expectedItemSaving,
		TaxClass:   expectedTaxClass,
	}, publicCart.Items[0])

	assert.Equal(t, db.OrderSummary{
		Total: db.TotalSummary{
			Value:  expectedTotal,
			Values: map[string]string{db.TaxClassTakeaway: expectedTotal},
		},
		Tax: db.TotalSummary{
			Value:  expectedTax,
			Values: map[string]string{db.TaxClassTakeaway: expectedTax},
		},
		Net: db.TotalSummary{
			Value:  expectedNet,
			Values: map[string]string{db.TaxClassTakeaway: expectedNet},
		},
		Saving: expectedCartSaving,
	}, publicCart.Summary)
}

func TestCalculate_TakeLastPayment(t *testing.T) {
	id := primitive.NewObjectID()
	sku := "some-sku"
	name := "product name"
	price := "12.34"
	expectedUnitPrice := "9.87"
	amount := 4
	expectedTotal := "39.48"
	expectedTax := "2.60"
	expectedNet := "36.88"
	expectedTaxClass := "takeaway"
	cartItem := db.CartItem{
		ProductID: id.Hex(),
		Options:   nil,
		Amount:    amount,
	}
	expectedItemSaving := "9.88"
	expectedCartSaving := "9.88"

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
	p1 := db.Payment{
		Session: checkout.CreateCheckoutSessionResponse{
			ExpiresAt: time.Now().Add(time.Hour),
			Amount: checkout.Amount{
				Currency: "EUR",
				Value:    3948,
			},
		},
		Environment: "env1",
		Client:      "client1",
	}
	p2 := db.Payment{
		Session: checkout.CreateCheckoutSessionResponse{
			ExpiresAt: time.Now().Add(time.Hour),
			Amount: checkout.Amount{
				Currency: "EUR",
				Value:    3948,
			},
		},
		Environment: "env2",
		Client:      "client2",
	}
	shoppingCart := db.Cart{
		ID:       id,
		IsPickup: true,
		Items:    []db.CartItem{cartItem},
		Payments: []db.Payment{p1, p2},
	}
	publicCart, err := cart.Calculate(context.TODO(), &shoppingCart, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, true, publicCart.IsPickup)
	assert.Equal(t, &p2, publicCart.Payment)

	assert.Equal(t, db.OrderItem{
		CartItem:   cartItem,
		SKU:        sku,
		Name:       name,
		Categories: categories,
		UnitPrice:  expectedUnitPrice,
		Total:      expectedTotal,
		Tax:        expectedTax,
		Net:        expectedNet,
		Saving:     expectedItemSaving,
		TaxClass:   expectedTaxClass,
	}, publicCart.Items[0])

	assert.Equal(t, db.OrderSummary{
		Total: db.TotalSummary{
			Value:  expectedTotal,
			Values: map[string]string{db.TaxClassTakeaway: expectedTotal},
		},
		Tax: db.TotalSummary{
			Value:  expectedTax,
			Values: map[string]string{db.TaxClassTakeaway: expectedTax},
		},
		Net: db.TotalSummary{
			Value:  expectedNet,
			Values: map[string]string{db.TaxClassTakeaway: expectedNet},
		},
		Saving: expectedCartSaving,
	}, publicCart.Summary)
}

func TestCalculate_ShippingMethodPickup(t *testing.T) {
	id := primitive.NewObjectID()
	sku := "some-sku"
	name := "product name"
	price := "12.34"
	expectedUnitPrice := "9.87"
	amount := 4
	expectedTotal := "39.48"
	expectedTax := "2.60"
	expectedNet := "36.88"
	expectedTaxClass := "takeaway"
	cartItem := db.CartItem{
		ProductID: id.Hex(),
		Options:   nil,
		Amount:    amount,
	}
	expectedItemSaving := "9.88"
	expectedCartSaving := "9.88"

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
		ID:       id,
		IsPickup: true,
		Items:    []db.CartItem{cartItem},
	}
	publicCart, err := cart.Calculate(context.TODO(), &shoppingCart, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, true, publicCart.IsPickup)

	assert.Equal(t, db.OrderItem{
		CartItem:   cartItem,
		SKU:        sku,
		Name:       name,
		Categories: categories,
		UnitPrice:  expectedUnitPrice,
		Total:      expectedTotal,
		Tax:        expectedTax,
		Net:        expectedNet,
		Saving:     expectedItemSaving,
		TaxClass:   expectedTaxClass,
	}, publicCart.Items[0])

	assert.Equal(t, db.OrderSummary{
		Total: db.TotalSummary{
			Value:  expectedTotal,
			Values: map[string]string{db.TaxClassTakeaway: expectedTotal},
		},
		Tax: db.TotalSummary{
			Value:  expectedTax,
			Values: map[string]string{db.TaxClassTakeaway: expectedTax},
		},
		Net: db.TotalSummary{
			Value:  expectedNet,
			Values: map[string]string{db.TaxClassTakeaway: expectedNet},
		},
		Saving: expectedCartSaving,
	}, publicCart.Summary)
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
		ProductID: id.Hex(),
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
		Paid:      true,
		Secret:    secret,
		CreatedAt: created,
		UpdatedAt: updated,
		Items:     []db.CartItem{cartItem},
		Payments:  nil,
	}
	publicCart, err := cart.Calculate(context.TODO(), &shoppingCart, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, id, publicCart.ID)
	assert.Equal(t, secret, publicCart.Secret)
	assert.Equal(t, updated, publicCart.UpdatedAt)
	assert.Equal(t, created, publicCart.CreatedAt)
	assert.Equal(t, false, publicCart.IsPickup)
	assert.True(t, publicCart.Paid)
	assert.Nil(t, publicCart.Payment)

	assert.Equal(t, db.OrderItem{
		CartItem:   cartItem,
		SKU:        sku,
		Name:       name,
		Categories: categories,
		UnitPrice:  price,
		Total:      expectedTotal,
		Tax:        expectedTax,
		Net:        expectedNet,
		TaxClass:   expectedTaxClass,
		Saving:     "0.00",
	}, publicCart.Items[0])

	assert.Equal(t, db.OrderSummary{
		Total: db.TotalSummary{
			Value:  expectedTotal,
			Values: map[string]string{db.TaxClassTakeaway: expectedTotal},
		},
		Tax: db.TotalSummary{
			Value:  expectedTax,
			Values: map[string]string{db.TaxClassTakeaway: expectedTax},
		},
		Net: db.TotalSummary{
			Value:  expectedNet,
			Values: map[string]string{db.TaxClassTakeaway: expectedNet},
		},
		Saving: "0.00",
	}, publicCart.Summary)
}

func TestCalculatePublicCart(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	colCarts, err := db.CollectionCarts(ctx)
	assert.NoError(t, err)
	cartID, err := primitive.ObjectIDFromHex("63c7a67ec0792b6ae57f57e7")
	assert.NoError(t, err)
	find := colCarts.FindOne(ctx, bson.M{"_id": cartID})
	shoppingCart := new(db.Cart)
	err = find.Decode(shoppingCart)
	assert.NoError(t, err)

	publicCart, err := cart.CalculatePublicCart(ctx, shoppingCart)
	assert.NoError(t, err)
	log.Printf("%+v", publicCart)
}
