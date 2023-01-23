package db

import (
	"context"
	"time"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/mongodb"
)

const colCoupons = "coupons"

func CollectionCoupons(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/shop/mongo/database", false)
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colCoupons)
}

const (
	CouponTypeGiftCard = "GiftCard"
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
