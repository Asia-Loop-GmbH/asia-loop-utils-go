package cart

import (
	"context"
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v9/pkg/shop/db"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func Test_generateCouponItems_NilCoupon(t *testing.T) {
	items := generateCouponItems(context.Background(), nil, &summaryTotal{
		Total:  decimal.Decimal{},
		Tax:    decimal.Decimal{},
		Net:    decimal.Decimal{},
		Saving: decimal.Decimal{},
	}, &summaryTotal{
		Total:  decimal.Decimal{},
		Tax:    decimal.Decimal{},
		Net:    decimal.Decimal{},
		Saving: decimal.Decimal{},
	})
	assert.Nil(t, items)
}

func Test_generateCouponItems_EinzweckCoupon(t *testing.T) {
	items := generateCouponItems(context.Background(), &db.Coupon{
		ID:        primitive.NewObjectID(),
		Type:      db.CouponTypeGiftCard,
		Code:      "GS1234",
		Total:     "10.00",
		Usage:     nil,
		Disabled:  false,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}, &summaryTotal{
		Total:  decimal.RequireFromString("12.00"),
		Tax:    decimal.RequireFromString("0.79"),
		Net:    decimal.RequireFromString("11.21"),
		Saving: decimal.Decimal{},
	}, &summaryTotal{
		Total:  decimal.Decimal{},
		Tax:    decimal.Decimal{},
		Net:    decimal.Decimal{},
		Saving: decimal.Decimal{},
	})
	assert.Equal(t, 1, len(items))
	assert.Equal(t, db.OrderItem{
		CartItem: db.CartItem{
			Amount: 1,
		},
		SKU:          db.CouponSKU,
		Name:         "Gutschein GS1234",
		Categories:   nil,
		UnitPrice:    "-10.00",
		Total:        "-10.00",
		Tax:          "-0.65",
		Net:          "-9.35",
		Saving:       "0.00",
		TaxClass:     db.TaxClassTakeaway,
		IsGiftCard:   false,
		GiftCardCode: nil,
	}, items[0])
}

func Test_generateCouponItems_EinzweckCoupon_CouponMoreThanCart(t *testing.T) {
	items := generateCouponItems(context.Background(), &db.Coupon{
		ID:        primitive.NewObjectID(),
		Type:      db.CouponTypeGiftCard,
		Code:      "GS1234",
		Total:     "10.00",
		Usage:     nil,
		Disabled:  false,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}, &summaryTotal{
		Total:  decimal.RequireFromString("8.00"),
		Tax:    decimal.RequireFromString("0.52"),
		Net:    decimal.RequireFromString("7.48"),
		Saving: decimal.Decimal{},
	}, &summaryTotal{
		Total:  decimal.RequireFromString("1.00"),
		Tax:    decimal.RequireFromString("0.16"),
		Net:    decimal.RequireFromString("0.84"),
		Saving: decimal.Decimal{},
	})
	assert.Equal(t, 1, len(items))
	assert.Equal(t, db.OrderItem{
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
	}, items[0])
}

func Test_generateCouponItems_MehrzweckCoupon_OnlyMitnehmen(t *testing.T) {
	items := generateCouponItems(context.Background(), &db.Coupon{
		ID:        primitive.NewObjectID(),
		Type:      db.CouponTypeMehrzweck,
		Code:      "GS1234",
		Total:     "5.00",
		Usage:     nil,
		Disabled:  false,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}, &summaryTotal{
		Total:  decimal.RequireFromString("8.00"),
		Tax:    decimal.RequireFromString("0.52"),
		Net:    decimal.RequireFromString("7.48"),
		Saving: decimal.Decimal{},
	}, &summaryTotal{
		Total:  decimal.RequireFromString("1.00"),
		Tax:    decimal.RequireFromString("0.16"),
		Net:    decimal.RequireFromString("0.84"),
		Saving: decimal.Decimal{},
	})
	assert.Equal(t, 1, len(items))
	assert.Equal(t, db.OrderItem{
		CartItem: db.CartItem{
			Amount: 1,
		},
		SKU:          db.CouponSKU,
		Name:         "Gutschein GS1234",
		Categories:   nil,
		UnitPrice:    "-5.00",
		Total:        "-5.00",
		Tax:          "-0.33",
		Net:          "-4.67",
		Saving:       "0.00",
		TaxClass:     db.TaxClassTakeaway,
		IsGiftCard:   false,
		GiftCardCode: nil,
	}, items[0])
}

func Test_generateCouponItems_MehrzweckCoupon_MitnehmenAndStandard(t *testing.T) {
	items := generateCouponItems(context.Background(), &db.Coupon{
		ID:        primitive.NewObjectID(),
		Type:      db.CouponTypeMehrzweck,
		Code:      "GS1234",
		Total:     "9.00",
		Usage:     nil,
		Disabled:  false,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}, &summaryTotal{
		Total:  decimal.RequireFromString("8.00"),
		Tax:    decimal.RequireFromString("0.52"),
		Net:    decimal.RequireFromString("7.48"),
		Saving: decimal.Decimal{},
	}, &summaryTotal{
		Total:  decimal.RequireFromString("5.00"),
		Tax:    decimal.RequireFromString("0.33"),
		Net:    decimal.RequireFromString("4.67"),
		Saving: decimal.Decimal{},
	})
	assert.Equal(t, 2, len(items))
	assert.Equal(t, db.OrderItem{
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
	}, items[0])
	assert.Equal(t, db.OrderItem{
		CartItem: db.CartItem{
			Amount: 1,
		},
		SKU:          db.CouponSKU,
		Name:         "Gutschein GS1234",
		Categories:   nil,
		UnitPrice:    "-1.00",
		Total:        "-1.00",
		Tax:          "-0.16",
		Net:          "-0.84",
		Saving:       "0.00",
		TaxClass:     db.TaxClassStandard,
		IsGiftCard:   false,
		GiftCardCode: nil,
	}, items[1])
}

func Test_generateCouponItems_MehrzweckCoupon_OnlyStandard(t *testing.T) {
	items := generateCouponItems(context.Background(), &db.Coupon{
		ID:        primitive.NewObjectID(),
		Type:      db.CouponTypeMehrzweck,
		Code:      "GS1234",
		Total:     "4.00",
		Usage:     nil,
		Disabled:  false,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}, &summaryTotal{
		Total:  decimal.RequireFromString("0.00"),
		Tax:    decimal.RequireFromString("0.00"),
		Net:    decimal.RequireFromString("0.00"),
		Saving: decimal.Decimal{},
	}, &summaryTotal{
		Total:  decimal.RequireFromString("5.00"),
		Tax:    decimal.RequireFromString("0.33"),
		Net:    decimal.RequireFromString("4.67"),
		Saving: decimal.Decimal{},
	})
	assert.Equal(t, 1, len(items))
	assert.Equal(t, db.OrderItem{
		CartItem: db.CartItem{
			Amount: 1,
		},
		SKU:          db.CouponSKU,
		Name:         "Gutschein GS1234",
		Categories:   nil,
		UnitPrice:    "-4.00",
		Total:        "-4.00",
		Tax:          "-0.64",
		Net:          "-3.36",
		Saving:       "0.00",
		TaxClass:     db.TaxClassStandard,
		IsGiftCard:   false,
		GiftCardCode: nil,
	}, items[0])
}

func Test_generateCouponItems_MehrzweckCoupon_CouponMoreThanCart(t *testing.T) {
	items := generateCouponItems(context.Background(), &db.Coupon{
		ID:        primitive.NewObjectID(),
		Type:      db.CouponTypeMehrzweck,
		Code:      "GS1234",
		Total:     "10.00",
		Usage:     nil,
		Disabled:  false,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}, &summaryTotal{
		Total:  decimal.RequireFromString("8.00"),
		Tax:    decimal.RequireFromString("0.52"),
		Net:    decimal.RequireFromString("7.48"),
		Saving: decimal.Decimal{},
	}, &summaryTotal{
		Total:  decimal.RequireFromString("1.00"),
		Tax:    decimal.RequireFromString("0.16"),
		Net:    decimal.RequireFromString("0.84"),
		Saving: decimal.Decimal{},
	})
	assert.Equal(t, 2, len(items))
	assert.Equal(t, db.OrderItem{
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
	}, items[0])
	assert.Equal(t, db.OrderItem{
		CartItem: db.CartItem{
			Amount: 1,
		},
		SKU:          db.CouponSKU,
		Name:         "Gutschein GS1234",
		Categories:   nil,
		UnitPrice:    "-1.00",
		Total:        "-1.00",
		Tax:          "-0.16",
		Net:          "-0.84",
		Saving:       "0.00",
		TaxClass:     db.TaxClassStandard,
		IsGiftCard:   false,
		GiftCardCode: nil,
	}, items[1])
}
