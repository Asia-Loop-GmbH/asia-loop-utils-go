package cart

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v6/pkg/orderutils"
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v6/pkg/shop/db"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/random"
)

func CreateOrder(ctx context.Context, shoppingCart *db.Cart, confirmedTotal string) (*db.Order, error) {
	log := logger.FromContext(ctx)
	log.Infof("Cart was paid, will create final order")

	colCoupons, err := db.CollectionCoupons(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init db collection")
	}
	colOrders, err := db.CollectionOrders(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init db collection")
	}

	invoiceNumber, err := orderutils.NextShopOrderInvoice(ctx)
	if err != nil {
		log.Errorf("Failed to generate invoice number: %s", err)
		return nil, errors.Wrap(err, "failed to generate invoice number")
	}
	orderNumber := fmt.Sprintf("S%s", random.String(10, lo.NumbersCharset))
	order, err := ToOrder(ctx, shoppingCart)
	if err != nil {
		log.Errorf("Failed to convert cart to order: %s", err)
		return nil, errors.Wrap(err, "failed to convert cart to order")
	}

	checkTotal := decimal.RequireFromString(confirmedTotal)
	orderTotal := decimal.RequireFromString(order.Summary.Total.Value)
	if !checkTotal.Equal(orderTotal) {
		log.Errorf("COUPON FRAUD DETECTED!!!: %s", orderNumber)
	}

	order.InvoiceNumber = invoiceNumber
	order.OrderNumber = &orderNumber
	for i := 0; i < len(order.Items); i++ {
		// check if gift card, then create codes
		if order.Items[i].IsGiftCard {
			order.Items[i].GiftCardCode = lo.Times(order.Items[i].Amount, func(index int) string {
				code := newGiftCardCode()
				log.Infof("Gift card code generated: %s = %s€", code, order.Items[i].UnitPrice)
				now := time.Now()
				coupon := db.Coupon{
					ID:        primitive.NewObjectID(),
					Type:      db.CouponTypeGiftCard,
					Code:      code,
					Total:     order.Items[i].UnitPrice,
					Usage:     make([]db.CouponUsage, 0),
					Disabled:  false,
					CreatedAt: now,
					UpdatedAt: now,
				}
				_, err := colCoupons.InsertOne(ctx, coupon)
				if err != nil {
					log.Errorf("Failed to save coupon [%s]", code)
					// TODO: we don't expect this happens, so it's good for now
				}
				return code
			})
		}

		// check if coupon, then update coupon value
		if order.Items[i].SKU == db.CouponSKU {
			findCoupon := colCoupons.FindOne(ctx, bson.M{"code": *order.CouponCode})
			usedCoupon := new(db.Coupon)
			err = findCoupon.Decode(usedCoupon)
			if err != nil {
				log.Errorf("Failed to find coupon [%s] to update: %s", *order.CouponCode, err)
			} else {
				_, err := colCoupons.UpdateByID(ctx, usedCoupon.ID, bson.D{{
					"$push", bson.D{
						{
							"usage",
							db.CouponUsage{
								OrderID:   order.ID.Hex(),
								Total:     decimal.RequireFromString(order.Items[i].Total).Neg().StringFixed(2),
								CreatedAt: time.Now(),
							},
						},
					},
				}})
				if err != nil {
					log.Errorf("Failed to update coupon [%s]: %s", *order.CouponCode, err)
				}
			}
		}
	}
	_, err = colOrders.InsertOne(ctx, order)
	if err != nil {
		log.Errorf("Failed to insert order: %s", err)
		return nil, errors.Wrap(err, "failed to insert order")
	}
	log.Infof("Order was created for cart [%s]", shoppingCart.ID)
	return order, nil
}

func newGiftCardCode() string {
	return fmt.Sprintf(
		"%s-%s-%s",
		random.String(4, lo.UpperCaseLettersCharset),
		random.String(4, lo.UpperCaseLettersCharset),
		random.String(4, lo.UpperCaseLettersCharset),
	)
}
