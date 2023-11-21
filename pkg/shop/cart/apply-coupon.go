package cart

import (
	"context"
	"fmt"
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v9/pkg/shop/db"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func generateCouponItems(ctx context.Context, coupon *db.Coupon, takeawaySummary *summaryTotal, standardSummary *summaryTotal) []db.OrderItem {
	log := logger.FromContext(ctx)
	log.Infof("Process coupon %v", coupon)

	if coupon == nil {
		log.Infof("Coupon is nil")
		return nil
	}

	switch coupon.Type {
	case db.CouponTypeGiftCard:
		return generateCouponItemsForEinzweck(ctx, coupon, takeawaySummary)
	case db.CouponTypeMehrzweck:
		return generateCouponItemsForMehrzweck(log, coupon, takeawaySummary, standardSummary)
	default:
		log.Infof("Unsupported coupon type [%v]", coupon.Type)
		return nil
	}
}

func generateCouponItemsForMehrzweck(log *logrus.Entry, coupon *db.Coupon, takeawaySummary *summaryTotal, standardSummary *summaryTotal) []db.OrderItem {
	log.Infof("Coupon is Mehrzweck")
	couponAmount := decimal.RequireFromString(coupon.Available())
	couponAppliedToMitnehmen := decimal.Min(couponAmount, takeawaySummary.Total)
	couponAmount = couponAmount.Sub(couponAppliedToMitnehmen)
	couponAppliedToStandard := decimal.Min(couponAmount, standardSummary.Total)
	couponMitnehmenTax := couponAppliedToMitnehmen.Div(decimal.NewFromFloat(1.07)).Mul(decimal.NewFromFloat(0.07)).Round(2)
	couponMitnehmenNet := couponAppliedToMitnehmen.Sub(couponMitnehmenTax)
	couponStandardTax := couponAppliedToStandard.Div(decimal.NewFromFloat(1.19)).Mul(decimal.NewFromFloat(0.19)).Round(2)
	couponStandardNet := couponAppliedToStandard.Sub(couponStandardTax)

	items := make([]db.OrderItem, 0)
	if !couponAppliedToMitnehmen.IsZero() {
		items = append(items, db.OrderItem{
			CartItem: db.CartItem{
				Amount: 1,
			},
			SKU:          db.CouponSKU,
			Name:         fmt.Sprintf("Gutschein %s", coupon.Code),
			Categories:   nil,
			UnitPrice:    couponAppliedToMitnehmen.Neg().StringFixed(2),
			Total:        couponAppliedToMitnehmen.Neg().StringFixed(2),
			Tax:          couponMitnehmenTax.Neg().StringFixed(2),
			Net:          couponMitnehmenNet.Neg().StringFixed(2),
			Saving:       "0.00",
			TaxClass:     db.TaxClassTakeaway,
			IsGiftCard:   false,
			GiftCardCode: nil,
		})
	}
	if !couponAppliedToStandard.IsZero() {
		items = append(items, db.OrderItem{
			CartItem: db.CartItem{
				Amount: 1,
			},
			SKU:          db.CouponSKU,
			Name:         fmt.Sprintf("Gutschein %s", coupon.Code),
			Categories:   nil,
			UnitPrice:    couponAppliedToStandard.Neg().StringFixed(2),
			Total:        couponAppliedToStandard.Neg().StringFixed(2),
			Tax:          couponStandardTax.Neg().StringFixed(2),
			Net:          couponStandardNet.Neg().StringFixed(2),
			Saving:       "0.00",
			TaxClass:     db.TaxClassStandard,
			IsGiftCard:   false,
			GiftCardCode: nil,
		})
	}

	return items
}

func generateCouponItemsForEinzweck(ctx context.Context, coupon *db.Coupon, takeawaySummary *summaryTotal) []db.OrderItem {
	log := logger.FromContext(ctx)
	log.Infof("Coupon is Einzweck")
	couponAmount := decimal.RequireFromString(coupon.Available())
	couponAmount = decimal.Min(couponAmount, takeawaySummary.Total)
	couponTax := couponAmount.Div(decimal.NewFromFloat(1.07)).Mul(decimal.NewFromFloat(0.07)).Round(2)
	couponNet := couponAmount.Sub(couponTax)

	return []db.OrderItem{
		{
			CartItem: db.CartItem{
				Amount: 1,
			},
			SKU:          db.CouponSKU,
			Name:         fmt.Sprintf("Gutschein %s", coupon.Code),
			Categories:   nil,
			UnitPrice:    couponAmount.Neg().StringFixed(2),
			Total:        couponAmount.Neg().StringFixed(2),
			Tax:          couponTax.Neg().StringFixed(2),
			Net:          couponNet.Neg().StringFixed(2),
			Saving:       "0.00",
			TaxClass:     db.TaxClassTakeaway,
			IsGiftCard:   false,
			GiftCardCode: nil,
		},
	}
}
