package db

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mongodb"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/random"
)

const colCoupons = "coupons"

func CollectionCoupons(ctx context.Context) (*mongo.Collection, error) {
	database, err := secretsmanager.GetParameter(ctx, "/shop/mongo/database")
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colCoupons)
}

const (
	CouponTypeGiftCard  = "GiftCard"  // This can be generated GS or old bought GS, a.k.a. Einzweck-Gutscheine
	CouponTypeMehrzweck = "Mehrzweck" // This can be the new bought GS, a.k.a. Mehrzweck-Gutscheine
)

type Coupon struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Type      string             `bson:"type" json:"type"`
	Code      string             `bson:"code" json:"code"`
	Total     string             `bson:"total" json:"total"`
	Usage     []CouponUsage      `bson:"usage" json:"usage"`
	Disabled  bool               `bson:"disabled" json:"disabled"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type CouponUsage struct {
	OrderID   string    `bson:"orderId" json:"orderId"`
	Total     string    `bson:"total" json:"total"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

func (c *Coupon) Available() string {
	if c.Disabled {
		return decimal.Zero.StringFixed(2)
	}

	available := decimal.RequireFromString(c.Total)
	lo.ForEach(c.Usage, func(u CouponUsage, _ int) {
		available = available.Sub(decimal.RequireFromString(u.Total))
	})
	return available.StringFixed(2)
}

func NewGiftCard(ctx context.Context, value string) (*Coupon, error) {
	log := logger.FromContext(ctx)
	log.Infof("Create new gift card")

	col, err := CollectionCoupons(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to init db collection")
	}

	code := newGifCardCode()
	log.Infof("Gift card code generated: %s = %s€", code, value)
	now := time.Now()
	coupon := Coupon{
		ID:        primitive.NewObjectID(),
		Type:      CouponTypeGiftCard,
		Code:      code,
		Total:     value,
		Usage:     make([]CouponUsage, 0),
		Disabled:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = col.InsertOne(ctx, coupon)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create coupon")
	}
	return &coupon, nil
}

func NewMehrzweckCoupon(ctx context.Context, value string) (*Coupon, error) {
	log := logger.FromContext(ctx)
	log.Infof("Create new Mehrzweck coupon")

	col, err := CollectionCoupons(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to init db collection")
	}

	code := newGifCardCode()
	log.Infof("Mehrzweck coupon code generated: %s = %s€", code, value)
	now := time.Now()
	coupon := Coupon{
		ID:        primitive.NewObjectID(),
		Type:      CouponTypeMehrzweck,
		Code:      code,
		Total:     value,
		Usage:     make([]CouponUsage, 0),
		Disabled:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = col.InsertOne(ctx, coupon)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create coupon")
	}
	return &coupon, nil
}

func newGifCardCode() string {
	return fmt.Sprintf(
		"%s-%s-%s",
		random.String(4, lo.UpperCaseLettersCharset),
		random.String(4, lo.UpperCaseLettersCharset),
		random.String(4, lo.UpperCaseLettersCharset),
	)
}
