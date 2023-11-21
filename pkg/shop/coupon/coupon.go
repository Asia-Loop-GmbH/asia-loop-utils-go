package coupon

import (
	"context"
	"fmt"
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v9/pkg/db"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"

	shopdb "github.com/asia-loop-gmbh/asia-loop-utils-go/v9/pkg/shop/db"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

func IsValidAndHasEnough(ctx context.Context, code, appliedAmount string) bool {
	log := logger.FromContext(ctx)
	log.Infof("check coupon valid and has enough amount: %s", code)

	coupon, err := GetCouponByCode(ctx, code)
	if err != nil {
		log.Errorf("could not get coupon '%s': %s", code, err)
		return false
	}
	current := decimal.RequireFromString(coupon.Available())
	toUse := decimal.RequireFromString(appliedAmount)
	return current.Cmp(toUse) >= 0
}

func GetCouponByCode(ctx context.Context, code string) (*shopdb.Coupon, error) {
	log := logger.FromContext(ctx)
	log.Infof("get coupon: %s", code)
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, fmt.Errorf("blank coupon code")
	}

	colCoupons, err := shopdb.CollectionCoupons(ctx)
	if err != nil {
		log.Errorf("Failed to init db collection")
		return nil, errors.Wrap(err, "failed to init db collection")
	}
	find := colCoupons.FindOne(ctx, bson.M{"code": strings.ToUpper(code)})
	shopCoupon := new(shopdb.Coupon)
	if err := find.Decode(shopCoupon); err != nil {
		log.Errorf("Failed to decode coupon: %s", err)
		return nil, errors.Wrap(err, "failed to decode coupon")
	}
	return shopCoupon, nil
}

var adminTaxClassToShopTaxClass = map[db.TaxClass]string{
	db.TaxClassStandard: shopdb.TaxClassStandard,
	db.TaxClassReduced:  shopdb.TaxClassTakeaway,
	db.TaxClassZero:     shopdb.TaxClassZero,
}

func UpdateCouponByAdminOrderItem(ctx context.Context, orderID string, item db.OrderItem) error {
	log := logger.FromContext(ctx)
	log.Infof("Update coupon by admin order item: %+v", item)
	code := strings.Split(item.Name, " ")[1]

	colCoupons, err := shopdb.CollectionCoupons(ctx)
	if err != nil {
		log.Errorf("Failed to init db collection")
		return errors.Wrap(err, "failed to init db collection")
	}

	update := colCoupons.FindOneAndUpdate(ctx, bson.M{"code": strings.ToUpper(code)}, bson.D{{
		"$push", bson.D{{"usage", shopdb.CouponUsage{
			OrderID:   orderID,
			Total:     decimal.RequireFromString(item.Total).Neg().StringFixed(2),
			Net:       lo.ToPtr(decimal.RequireFromString(item.Net).Neg().StringFixed(2)),
			Tax:       lo.ToPtr(decimal.RequireFromString(item.Tax).Neg().StringFixed(2)),
			TaxClass:  lo.ToPtr(adminTaxClassToShopTaxClass[item.TaxClass]),
			CreatedAt: time.Now(),
		}}},
	}})

	if err := update.Err(); err != nil {
		log.Errorf("Failed to update coupon: %s", err)
		return errors.Wrap(err, "failed to update coupon")
	}
	log.Infof("Shop coupon [%s] updated", code)
	return nil
}

func UpdateCouponByOrderItem(ctx context.Context, orderID string, item shopdb.OrderItem) error {
	log := logger.FromContext(ctx)
	log.Infof("Update coupon: %+v", item)
	code := strings.Split(item.Name, " ")[1]

	colCoupons, err := shopdb.CollectionCoupons(ctx)
	if err != nil {
		log.Errorf("Failed to init db collection")
		return errors.Wrap(err, "failed to init db collection")
	}

	update := colCoupons.FindOneAndUpdate(ctx, bson.M{"code": strings.ToUpper(code)}, bson.D{{
		"$push", bson.D{{"usage", shopdb.CouponUsage{
			OrderID:   orderID,
			Total:     decimal.RequireFromString(item.Total).Neg().StringFixed(2),
			Net:       lo.ToPtr(decimal.RequireFromString(item.Net).Neg().StringFixed(2)),
			Tax:       lo.ToPtr(decimal.RequireFromString(item.Tax).Neg().StringFixed(2)),
			TaxClass:  lo.ToPtr(item.TaxClass),
			CreatedAt: time.Now(),
		}}},
	}})

	if err := update.Err(); err != nil {
		log.Errorf("Failed to update coupon: %s", err)
		return errors.Wrap(err, "failed to update coupon")
	}
	log.Infof("Shop coupon [%s] updated", code)
	return nil
}

func UpdateCouponByCode(ctx context.Context, code, amount string) error {
	log := logger.FromContext(ctx)
	log.Infof("update coupon %s: %s", code, amount)

	toUse := decimal.RequireFromString(amount)

	colCoupons, err := shopdb.CollectionCoupons(ctx)
	if err != nil {
		log.Errorf("Failed to init db collection")
		return errors.Wrap(err, "failed to init db collection")
	}

	update := colCoupons.FindOneAndUpdate(ctx, bson.M{"code": strings.ToUpper(code)}, bson.D{{
		"$push", bson.D{{"usage", shopdb.CouponUsage{
			OrderID:   "Local in admin app",
			Total:     toUse.StringFixed(2),
			CreatedAt: time.Now(),
		}}},
	}})
	if err := update.Err(); err != nil {
		log.Errorf("Failed to update coupon: %s", err)
		return errors.Wrap(err, "failed to update coupon")
	}
	log.Infof("Shop coupon [%s] updated", code)
	return nil
}
