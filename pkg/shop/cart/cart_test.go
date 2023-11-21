package cart

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/adyen/adyen-go-api-library/v8/src/checkout"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v9/pkg/shop/db"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
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
		{
			ID:          primitive.NewObjectID(),
			Name:        "standard",
			DisplayName: "MwSt. 19%",
			Rate:        "0.19",
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "kein",
			DisplayName: "MwSt. 0%",
			Rate:        "0.00",
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

func TestToOrder_IgnoreExpired(t *testing.T) {
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
	order, err := toOrder(context.TODO(), &shoppingCart, nil, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, true, order.IsPickup)
	assert.Nil(t, order.Payment)

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
	}, order.Items[0])

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
	}, order.Summary)
}

func TestToOrder_IgnoreTotalNotMatch(t *testing.T) {
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
	order, err := toOrder(context.TODO(), &shoppingCart, nil, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, true, order.IsPickup)
	assert.Nil(t, order.Payment)

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
	}, order.Items[0])

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
	}, order.Summary)
}

func TestToOrder_TakeLastPayment(t *testing.T) {
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
	order, err := toOrder(context.TODO(), &shoppingCart, nil, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, true, order.IsPickup)
	assert.Equal(t, &p2, order.Payment)

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
	}, order.Items[0])

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
	}, order.Summary)
}

func TestToOrder_ShippingMethodPickup(t *testing.T) {
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
	order, err := toOrder(context.TODO(), &shoppingCart, nil, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, true, order.IsPickup)

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
	}, order.Items[0])

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
	}, order.Summary)
}

func TestToOrder_VariableProduct(t *testing.T) {
	id := primitive.NewObjectID()
	secret := "some-secret"
	updated := time.Now()
	created := time.Now()
	sku := "some-sku"
	name := "product name"
	price := "12.34"
	varPrice := "15.67"
	amount := 4
	expectedTotal := "62.68"
	expectedTax := "4.12"
	expectedNet := "58.56"
	expectedTaxClass := "takeaway"
	categories := []string{"category1", "category2"}
	variableProduct := db.Product{
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
		Options:          []string{"option-1", "option-2"},
		Variations: []db.ProductVariation{
			{
				SKU: sku,
				Price: db.CustomizablePrice{
					Value:        varPrice,
					CustomValues: nil,
				},
				Options: []db.ProductVariationOption{
					{
						Name:  "option-1",
						Value: "value-1",
					},
				},
			},
		},
	}
	cartItem := db.CartItem{
		ProductID: id.Hex(),
		Options: map[string]string{
			"option-1": "value-1",
			"option-2": "any-value",
		},
		Amount: amount,
	}
	products := []db.Product{variableProduct}
	storeKey := "ERLANGEN"
	cartCheckout := &db.CartCheckout{
		FirstName:    "First Name",
		LastName:     "Last",
		AddressLine1: "Line1",
		AddressLine2: "",
		City:         "City",
		Postcode:     "Postcode",
		Telephone:    "Tel",
		Email:        "Email",
		Note:         "",
		Date:         "",
		Slot:         "",
		Begin:        lo.ToPtr(time.Now()),
	}
	shoppingCart := db.Cart{
		ID:        id,
		StoreKey:  storeKey,
		IsPickup:  false,
		Paid:      true,
		Secret:    secret,
		CreatedAt: created,
		UpdatedAt: updated,
		Items:     []db.CartItem{cartItem},
		Payments:  nil,
		Checkout:  cartCheckout,
	}
	order, err := toOrder(context.TODO(), &shoppingCart, nil, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, id, order.ID)
	assert.Equal(t, secret, order.Secret)
	assert.Equal(t, updated, order.UpdatedAt)
	assert.Equal(t, created, order.CreatedAt)
	assert.Equal(t, false, order.IsPickup)
	assert.True(t, order.Paid)
	assert.Equal(t, storeKey, order.StoreKey)
	assert.Equal(t, cartCheckout, order.Checkout)
	assert.Nil(t, order.Payment)

	assert.Equal(t, db.OrderItem{
		CartItem:   cartItem,
		SKU:        sku,
		Name:       name,
		Categories: categories,
		UnitPrice:  varPrice,
		Total:      expectedTotal,
		Tax:        expectedTax,
		Net:        expectedNet,
		TaxClass:   expectedTaxClass,
		Saving:     "0.00",
	}, order.Items[0])

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
	}, order.Summary)
}

func TestToOrder_StorePrice(t *testing.T) {
	id := primitive.NewObjectID()
	secret := "some-secret"
	updated := time.Now()
	created := time.Now()
	sku := "some-sku"
	name := "product name"
	price := "12.34"
	varPrice := "15.67"
	amount := 4
	expectedTotal := "62.68"
	expectedTax := "4.12"
	expectedNet := "58.56"
	expectedTaxClass := "takeaway"
	categories := []string{"category1", "category2"}
	storeKey := "ERLANGEN"
	variableProduct := db.Product{
		ID:   id,
		SKU:  sku,
		Name: name,
		Price: db.CustomizablePrice{
			Value: price,
			CustomValues: map[string]string{
				storeKey: varPrice,
			},
		},
		TaxClassStandard: "standard",
		TaxClassTakeAway: "takeaway",
		Categories:       categories,
	}
	cartItem := db.CartItem{
		ProductID: id.Hex(),
		Amount:    amount,
	}
	products := []db.Product{variableProduct}
	cartCheckout := &db.CartCheckout{
		FirstName:    "First Name",
		LastName:     "Last",
		AddressLine1: "Line1",
		AddressLine2: "",
		City:         "City",
		Postcode:     "Postcode",
		Telephone:    "Tel",
		Email:        "Email",
		Note:         "",
		Date:         "",
		Slot:         "",
		Begin:        lo.ToPtr(time.Now()),
	}
	shoppingCart := db.Cart{
		ID:        id,
		StoreKey:  storeKey,
		IsPickup:  false,
		Paid:      true,
		Secret:    secret,
		CreatedAt: created,
		UpdatedAt: updated,
		Items:     []db.CartItem{cartItem},
		Payments:  nil,
		Checkout:  cartCheckout,
	}
	order, err := toOrder(context.TODO(), &shoppingCart, nil, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, id, order.ID)
	assert.Equal(t, secret, order.Secret)
	assert.Equal(t, updated, order.UpdatedAt)
	assert.Equal(t, created, order.CreatedAt)
	assert.Equal(t, false, order.IsPickup)
	assert.True(t, order.Paid)
	assert.Equal(t, storeKey, order.StoreKey)
	assert.Equal(t, cartCheckout, order.Checkout)
	assert.Nil(t, order.Payment)

	assert.Equal(t, db.OrderItem{
		CartItem:   cartItem,
		SKU:        sku,
		Name:       name,
		Categories: categories,
		UnitPrice:  varPrice,
		Total:      expectedTotal,
		Tax:        expectedTax,
		Net:        expectedNet,
		TaxClass:   expectedTaxClass,
		Saving:     "0.00",
	}, order.Items[0])

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
	}, order.Summary)
}

func TestToOrder_ApplyCoupon_CouponBiggerThanTotal(t *testing.T) {
	id := primitive.NewObjectID()
	secret := "some-secret"
	updated := time.Now()
	created := time.Now()
	sku := "some-sku"
	name := "product name"
	price := "12.34"
	amount := 1
	expectedTotal := "0.00"
	expectedTax := "0.00"
	expectedNet := "0.00"
	categories := []string{"category1", "category2"}
	storeKey := "ERLANGEN"
	variableProduct := db.Product{
		ID:   id,
		SKU:  sku,
		Name: name,
		Price: db.CustomizablePrice{
			Value: price,
		},
		TaxClassStandard: db.TaxClassStandard,
		TaxClassTakeAway: db.TaxClassTakeaway,
		Categories:       categories,
	}
	cartItem := db.CartItem{
		ProductID: id.Hex(),
		Amount:    amount,
	}
	coupon := &db.Coupon{
		Type:  db.CouponTypeGiftCard,
		Code:  "GS1234",
		Total: "20.00",
		Usage: []db.CouponUsage{
			{
				Total: "2.00",
			},
		},
		Disabled: false,
	}
	expectedOrderItem := db.OrderItem{
		CartItem: db.CartItem{
			Amount: 1,
		},
		SKU:          db.CouponSKU,
		Name:         "Gutschein GS1234",
		Categories:   nil,
		UnitPrice:    "-12.34",
		Total:        "-12.34",
		Tax:          "-0.81",
		Net:          "-11.53",
		Saving:       "0.00",
		TaxClass:     db.TaxClassTakeaway,
		IsGiftCard:   false,
		GiftCardCode: nil,
	}
	products := []db.Product{variableProduct}
	cartCheckout := &db.CartCheckout{
		FirstName:    "First Name",
		LastName:     "Last",
		AddressLine1: "Line1",
		AddressLine2: "",
		City:         "City",
		Postcode:     "Postcode",
		Telephone:    "Tel",
		Email:        "Email",
		Note:         "",
		Date:         "",
		Slot:         "",
		Begin:        lo.ToPtr(time.Now()),
	}
	shoppingCart := db.Cart{
		ID:        id,
		StoreKey:  storeKey,
		IsPickup:  false,
		Paid:      true,
		Secret:    secret,
		CreatedAt: created,
		UpdatedAt: updated,
		Items:     []db.CartItem{cartItem},
		Payments:  nil,
		Checkout:  cartCheckout,
	}
	order, err := toOrder(context.TODO(), &shoppingCart, coupon, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, id, order.ID)
	assert.Equal(t, secret, order.Secret)
	assert.Equal(t, updated, order.UpdatedAt)
	assert.Equal(t, created, order.CreatedAt)
	assert.Equal(t, false, order.IsPickup)
	assert.True(t, order.Paid)
	assert.Equal(t, storeKey, order.StoreKey)
	assert.Equal(t, cartCheckout, order.Checkout)
	assert.Nil(t, order.Payment)

	assert.Equal(t, 2, len(order.Items))
	// order.Items[0] is tested by other methods
	assert.Equal(t, expectedOrderItem, order.Items[1])

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
	}, order.Summary)
}

func TestToOrder_ApplyCoupon(t *testing.T) {
	id := primitive.NewObjectID()
	secret := "some-secret"
	updated := time.Now()
	created := time.Now()
	sku := "some-sku"
	name := "product name"
	price := "12.34"
	amount := 1
	expectedTotal := "4.34"
	expectedTax := "0.29"
	expectedNet := "4.05"
	categories := []string{"category1", "category2"}
	storeKey := "ERLANGEN"
	variableProduct := db.Product{
		ID:   id,
		SKU:  sku,
		Name: name,
		Price: db.CustomizablePrice{
			Value: price,
		},
		TaxClassStandard: db.TaxClassStandard,
		TaxClassTakeAway: db.TaxClassTakeaway,
		Categories:       categories,
	}
	cartItem := db.CartItem{
		ProductID: id.Hex(),
		Amount:    amount,
	}
	coupon := &db.Coupon{
		Type:  db.CouponTypeGiftCard,
		Total: "10.00",
		Code:  "GS1234",
		Usage: []db.CouponUsage{
			{
				Total: "2.00",
			},
		},
		Disabled: false,
	}
	expectedOrderItem := db.OrderItem{
		CartItem: db.CartItem{
			Amount: 1,
		},
		SKU:          db.CouponSKU,
		Name:         "Gutschein GS1234",
		Categories:   nil,
		UnitPrice:    "-8.00",
		Total:        "-8.00",
		Tax:          "-0.52",
		Net:          "-7.48",
		Saving:       "0.00",
		TaxClass:     db.TaxClassTakeaway,
		IsGiftCard:   false,
		GiftCardCode: nil,
	}
	products := []db.Product{variableProduct}
	cartCheckout := &db.CartCheckout{
		FirstName:    "First Name",
		LastName:     "Last",
		AddressLine1: "Line1",
		AddressLine2: "",
		City:         "City",
		Postcode:     "Postcode",
		Telephone:    "Tel",
		Email:        "Email",
		Note:         "",
		Date:         "",
		Slot:         "",
		Begin:        lo.ToPtr(time.Now()),
	}
	couponCode := lo.ToPtr("coupon")
	shoppingCart := db.Cart{
		ID:         id,
		StoreKey:   storeKey,
		CouponCode: couponCode,
		IsPickup:   false,
		Paid:       true,
		Secret:     secret,
		CreatedAt:  created,
		UpdatedAt:  updated,
		Items:      []db.CartItem{cartItem},
		Payments:   nil,
		Checkout:   cartCheckout,
	}
	order, err := toOrder(context.TODO(), &shoppingCart, coupon, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, id, order.ID)
	assert.Equal(t, secret, order.Secret)
	assert.Equal(t, updated, order.UpdatedAt)
	assert.Equal(t, created, order.CreatedAt)
	assert.Equal(t, false, order.IsPickup)
	assert.True(t, order.Paid)
	assert.Equal(t, storeKey, order.StoreKey)
	assert.Equal(t, cartCheckout, order.Checkout)
	assert.Nil(t, order.Payment)
	assert.Equal(t, couponCode, order.CouponCode)

	assert.Equal(t, 2, len(order.Items))
	// order.Items[0] is tested by other methods
	assert.Equal(t, expectedOrderItem, order.Items[1])

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
	}, order.Summary)
}

func TestToOrder_GiftCard(t *testing.T) {
	id := primitive.NewObjectID()
	secret := "some-secret"
	updated := time.Now()
	created := time.Now()
	sku := "some-sku"
	name := "product name"
	price := "50.00"
	expectedTotal := "50.00"
	expectedNet := "50.00"
	expectedTax := "0.00"
	expectedTaxClass := "kein"
	categories := []string{"category1", "category2"}
	storeKey := "ERLANGEN"
	giftCard := db.Product{
		ID:   id,
		SKU:  sku,
		Name: name,
		Price: db.CustomizablePrice{
			Value: price,
		},
		TaxClassStandard: "kein",
		TaxClassTakeAway: "kein",
		IsGiftCard:       true,
		Categories:       categories,
	}
	cartItem := db.CartItem{
		ProductID: id.Hex(),
		Amount:    1,
	}
	products := []db.Product{giftCard}
	cartCheckout := &db.CartCheckout{
		FirstName:    "First Name",
		LastName:     "Last",
		AddressLine1: "Line1",
		AddressLine2: "",
		City:         "City",
		Postcode:     "Postcode",
		Telephone:    "Tel",
		Email:        "Email",
		Note:         "",
		Date:         "",
		Slot:         "",
		Begin:        lo.ToPtr(time.Now()),
	}
	shoppingCart := db.Cart{
		ID:        id,
		StoreKey:  storeKey,
		IsPickup:  true,
		Paid:      true,
		Secret:    secret,
		CreatedAt: created,
		UpdatedAt: updated,
		Items:     []db.CartItem{cartItem},
		Payments:  nil,
		Checkout:  cartCheckout,
	}
	order, err := toOrder(context.TODO(), &shoppingCart, nil, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, id, order.ID)
	assert.Equal(t, secret, order.Secret)
	assert.Equal(t, updated, order.UpdatedAt)
	assert.Equal(t, created, order.CreatedAt)
	assert.Equal(t, true, order.IsPickup)
	assert.True(t, order.Paid)
	assert.Equal(t, storeKey, order.StoreKey)
	assert.Equal(t, cartCheckout, order.Checkout)
	assert.Nil(t, order.Payment)

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
		IsGiftCard: true,
		Saving:     "0.00",
	}, order.Items[0])

	assert.Equal(t, db.OrderSummary{
		Total: db.TotalSummary{
			Value:  expectedTotal,
			Values: map[string]string{db.TaxClassZero: expectedTotal, db.TaxClassTakeaway: "0.00"},
		},
		Tax: db.TotalSummary{
			Value:  expectedTax,
			Values: map[string]string{db.TaxClassTakeaway: "0.00"},
		},
		Net: db.TotalSummary{
			Value:  expectedNet,
			Values: map[string]string{db.TaxClassZero: expectedNet, db.TaxClassTakeaway: "0.00"},
		},
		Saving: "0.00",
	}, order.Summary)
}

func TestToOrder_IgnoreZeroCustomPrice(t *testing.T) {
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
	storeKey := "ERLANGEN"
	categories := []string{"category1", "category2"}
	simpleProduct := db.Product{
		ID:   id,
		SKU:  sku,
		Name: name,
		Price: db.CustomizablePrice{
			Value: price,
			CustomValues: map[string]string{
				storeKey: "0.00",
			},
		},
		TaxClassStandard: "standard",
		TaxClassTakeAway: "takeaway",
		Categories:       categories,
		Options:          nil,
		Variations:       nil,
	}
	products := []db.Product{simpleProduct}
	cartCheckout := &db.CartCheckout{
		FirstName:    "First Name",
		LastName:     "Last",
		AddressLine1: "Line1",
		AddressLine2: "",
		City:         "City",
		Postcode:     "Postcode",
		Telephone:    "Tel",
		Email:        "Email",
		Note:         "",
		Date:         "",
		Slot:         "",
		Begin:        lo.ToPtr(time.Now()),
	}
	shoppingCart := db.Cart{
		ID:        id,
		StoreKey:  storeKey,
		IsPickup:  false,
		Paid:      true,
		Secret:    secret,
		CreatedAt: created,
		UpdatedAt: updated,
		Items:     []db.CartItem{cartItem},
		Payments:  nil,
		Checkout:  cartCheckout,
	}
	order, err := toOrder(context.TODO(), &shoppingCart, nil, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, id, order.ID)
	assert.Equal(t, secret, order.Secret)
	assert.Equal(t, updated, order.UpdatedAt)
	assert.Equal(t, created, order.CreatedAt)
	assert.Equal(t, false, order.IsPickup)
	assert.True(t, order.Paid)
	assert.Equal(t, storeKey, order.StoreKey)
	assert.Equal(t, cartCheckout, order.Checkout)
	assert.Nil(t, order.Payment)

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
	}, order.Items[0])

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
	}, order.Summary)
}

func TestToOrder_Tip(t *testing.T) {
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
	simpleProduct := db.Product{
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
	}
	products := []db.Product{simpleProduct}
	storeKey := "ERLANGEN"
	cartCheckout := &db.CartCheckout{
		FirstName:    "First Name",
		LastName:     "Last",
		AddressLine1: "Line1",
		AddressLine2: "",
		City:         "City",
		Postcode:     "Postcode",
		Telephone:    "Tel",
		Email:        "Email",
		Note:         "",
		Date:         "",
		Slot:         "",
		Begin:        lo.ToPtr(time.Now()),
	}
	user := lo.ToPtr("some user")
	tip := "1.23"
	shoppingCart := db.Cart{
		ID:        id,
		User:      user,
		StoreKey:  storeKey,
		Tip:       lo.ToPtr(tip),
		IsPickup:  false,
		Paid:      true,
		Secret:    secret,
		CreatedAt: created,
		UpdatedAt: updated,
		Items:     []db.CartItem{cartItem},
		Payments:  nil,
		Checkout:  cartCheckout,
	}
	order, err := toOrder(context.TODO(), &shoppingCart, nil, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, id, order.ID)
	assert.Equal(t, secret, order.Secret)
	assert.Equal(t, updated, order.UpdatedAt)
	assert.Equal(t, created, order.CreatedAt)
	assert.Equal(t, false, order.IsPickup)
	assert.True(t, order.Paid)
	assert.Equal(t, storeKey, order.StoreKey)
	assert.Equal(t, cartCheckout, order.Checkout)
	assert.Equal(t, user, order.User)
	assert.Equal(t, tip, *order.Tip)
	assert.Nil(t, order.Payment)

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
	}, order.Items[0])
	assert.Equal(t, db.OrderItem{
		CartItem: db.CartItem{
			Amount: 1,
		},
		SKU:          db.TipSKU,
		Name:         "Trinkgeld",
		Categories:   nil,
		UnitPrice:    tip,
		Total:        tip,
		Tax:          "0.00",
		Net:          tip,
		Saving:       "0.00",
		TaxClass:     db.TaxClassZero,
		IsGiftCard:   false,
		GiftCardCode: nil,
	}, order.Items[1])

	expectedTotalWithTip := decimal.RequireFromString(expectedTotal).Add(decimal.RequireFromString(tip)).StringFixed(2)
	expectedNetWithTip := decimal.RequireFromString(expectedNet).Add(decimal.RequireFromString(tip)).StringFixed(2)
	assert.Equal(t, db.OrderSummary{
		Total: db.TotalSummary{
			Value: expectedTotalWithTip,
			Values: map[string]string{
				db.TaxClassTakeaway: expectedTotal,
				db.TaxClassZero:     tip,
			},
		},
		Tax: db.TotalSummary{
			Value:  expectedTax,
			Values: map[string]string{db.TaxClassTakeaway: expectedTax},
		},
		Net: db.TotalSummary{
			Value: expectedNetWithTip,
			Values: map[string]string{
				db.TaxClassTakeaway: expectedNet,
				db.TaxClassZero:     tip,
			},
		},
		Saving: "0.00",
	}, order.Summary)
}

func TestToOrder_Tip_And_GiftCard(t *testing.T) {
	id := primitive.NewObjectID()
	secret := "some-secret"
	updated := time.Now()
	created := time.Now()
	tip := "10.00"
	sku := "some-sku"
	name := "product name"
	price := "50.00"
	expectedTotal := "60.00"
	expectedNet := "60.00"
	expectedTax := "0.00"
	expectedTaxClass := "kein"
	categories := []string{"category1", "category2"}
	storeKey := "ERLANGEN"
	giftCard := db.Product{
		ID:   id,
		SKU:  sku,
		Name: name,
		Price: db.CustomizablePrice{
			Value: price,
		},
		TaxClassStandard: "kein",
		TaxClassTakeAway: "kein",
		IsGiftCard:       true,
		Categories:       categories,
	}
	cartItem := db.CartItem{
		ProductID: id.Hex(),
		Amount:    1,
	}
	products := []db.Product{giftCard}
	cartCheckout := &db.CartCheckout{
		FirstName:    "First Name",
		LastName:     "Last",
		AddressLine1: "Line1",
		AddressLine2: "",
		City:         "City",
		Postcode:     "Postcode",
		Telephone:    "Tel",
		Email:        "Email",
		Note:         "",
		Date:         "",
		Slot:         "",
		Begin:        lo.ToPtr(time.Now()),
	}
	shoppingCart := db.Cart{
		ID:        id,
		StoreKey:  storeKey,
		IsPickup:  true,
		Paid:      true,
		Secret:    secret,
		CreatedAt: created,
		UpdatedAt: updated,
		Items:     []db.CartItem{cartItem},
		Payments:  nil,
		Checkout:  cartCheckout,
		Tip:       lo.ToPtr(tip),
	}
	order, err := toOrder(context.TODO(), &shoppingCart, nil, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, id, order.ID)
	assert.Equal(t, secret, order.Secret)
	assert.Equal(t, updated, order.UpdatedAt)
	assert.Equal(t, created, order.CreatedAt)
	assert.Equal(t, true, order.IsPickup)
	assert.True(t, order.Paid)
	assert.Equal(t, storeKey, order.StoreKey)
	assert.Equal(t, cartCheckout, order.Checkout)
	assert.Nil(t, order.Payment)

	assert.Equal(t, db.OrderItem{
		CartItem:   cartItem,
		SKU:        sku,
		Name:       name,
		Categories: categories,
		UnitPrice:  price,
		Total:      price,
		Tax:        "0.00",
		Net:        price,
		TaxClass:   expectedTaxClass,
		IsGiftCard: true,
		Saving:     "0.00",
	}, order.Items[0])

	assert.Equal(t, db.OrderSummary{
		Total: db.TotalSummary{
			Value:  expectedTotal,
			Values: map[string]string{db.TaxClassZero: expectedTotal, db.TaxClassTakeaway: "0.00"},
		},
		Tax: db.TotalSummary{
			Value:  expectedTax,
			Values: map[string]string{db.TaxClassTakeaway: "0.00"},
		},
		Net: db.TotalSummary{
			Value:  expectedNet,
			Values: map[string]string{db.TaxClassZero: expectedNet, db.TaxClassTakeaway: "0.00"},
		},
		Saving: "0.00",
	}, order.Summary)
}

func TestToOrder_Drink(t *testing.T) {
	id := primitive.NewObjectID()
	secret := "some-secret"
	updated := time.Now()
	created := time.Now()
	sku := "some-sku"
	name := "product name"
	price := "12.34" // net 10.37, tax 1.97
	amount := 4
	expectedTotal := "49.36"
	expectedTax := "7.88"
	expectedNet := "41.48"
	expectedTaxClass := "standard"
	cartItem := db.CartItem{
		ProductID: id.Hex(),
		Options:   nil,
		Amount:    amount,
	}
	categories := []string{"category1", "category2", "drink"}
	simpleProduct := db.Product{
		ID:   id,
		SKU:  sku,
		Name: name,
		Price: db.CustomizablePrice{
			Value:        price,
			CustomValues: nil,
		},
		TaxClassStandard: "standard",
		TaxClassTakeAway: "standard",
		Categories:       categories,
		Options:          nil,
		Variations:       nil,
	}
	products := []db.Product{simpleProduct}
	storeKey := "ERLANGEN"
	cartCheckout := &db.CartCheckout{
		FirstName:    "First Name",
		LastName:     "Last",
		AddressLine1: "Line1",
		AddressLine2: "",
		City:         "City",
		Postcode:     "Postcode",
		Telephone:    "Tel",
		Email:        "Email",
		Note:         "",
		Date:         "",
		Slot:         "",
		Begin:        lo.ToPtr(time.Now()),
	}
	user := lo.ToPtr("some user")
	shoppingCart := db.Cart{
		ID:        id,
		User:      user,
		StoreKey:  storeKey,
		IsPickup:  true,
		Paid:      true,
		Secret:    secret,
		CreatedAt: created,
		UpdatedAt: updated,
		Items:     []db.CartItem{cartItem},
		Payments:  nil,
		Checkout:  cartCheckout,
	}
	order, err := toOrder(context.TODO(), &shoppingCart, nil, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, id, order.ID)
	assert.Equal(t, secret, order.Secret)
	assert.Equal(t, updated, order.UpdatedAt)
	assert.Equal(t, created, order.CreatedAt)
	assert.Equal(t, true, order.IsPickup)
	assert.True(t, order.Paid)
	assert.Equal(t, storeKey, order.StoreKey)
	assert.Equal(t, cartCheckout, order.Checkout)
	assert.Equal(t, user, order.User)
	assert.Nil(t, order.Payment)

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
	}, order.Items[0])

	assert.Equal(t, db.OrderSummary{
		Total: db.TotalSummary{
			Value:  expectedTotal,
			Values: map[string]string{db.TaxClassStandard: expectedTotal, db.TaxClassTakeaway: "0.00"},
		},
		Tax: db.TotalSummary{
			Value:  expectedTax,
			Values: map[string]string{db.TaxClassStandard: expectedTax, db.TaxClassTakeaway: "0.00"},
		},
		Net: db.TotalSummary{
			Value:  expectedNet,
			Values: map[string]string{db.TaxClassStandard: expectedNet, db.TaxClassTakeaway: "0.00"},
		},
		Saving: "0.00",
	}, order.Summary)
}

func TestToOrder(t *testing.T) {
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
	simpleProduct := db.Product{
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
	}
	products := []db.Product{simpleProduct}
	storeKey := "ERLANGEN"
	cartCheckout := &db.CartCheckout{
		FirstName:    "First Name",
		LastName:     "Last",
		AddressLine1: "Line1",
		AddressLine2: "",
		City:         "City",
		Postcode:     "Postcode",
		Telephone:    "Tel",
		Email:        "Email",
		Note:         "",
		Date:         "",
		Slot:         "",
		Begin:        lo.ToPtr(time.Now()),
	}
	user := lo.ToPtr("some user")
	shoppingCart := db.Cart{
		ID:        id,
		User:      user,
		StoreKey:  storeKey,
		IsPickup:  false,
		Paid:      true,
		Secret:    secret,
		CreatedAt: created,
		UpdatedAt: updated,
		Items:     []db.CartItem{cartItem},
		Payments:  nil,
		Checkout:  cartCheckout,
	}
	order, err := toOrder(context.TODO(), &shoppingCart, nil, products, taxes)

	assert.NoError(t, err)
	assert.Equal(t, id, order.ID)
	assert.Equal(t, secret, order.Secret)
	assert.Equal(t, updated, order.UpdatedAt)
	assert.Equal(t, created, order.CreatedAt)
	assert.Equal(t, false, order.IsPickup)
	assert.True(t, order.Paid)
	assert.Equal(t, storeKey, order.StoreKey)
	assert.Equal(t, cartCheckout, order.Checkout)
	assert.Equal(t, user, order.User)
	assert.Nil(t, order.Payment)

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
	}, order.Items[0])

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
	}, order.Summary)
}

func TestToOrder_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	colCarts, err := db.CollectionCarts(ctx)
	assert.NoError(t, err)
	cartID, err := primitive.ObjectIDFromHex("63ceb8f5f6ec033180dec5c1")
	assert.NoError(t, err)
	find := colCarts.FindOne(ctx, bson.M{"_id": cartID})
	shoppingCart := new(db.Cart)
	err = find.Decode(shoppingCart)
	assert.NoError(t, err)

	order, err := ToOrder(ctx, shoppingCart)
	assert.NoError(t, err)
	log.Printf("%+v", order)
}
