package cart

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v7/pkg/orderutils"
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v7/pkg/shop/db"
	mysns "github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/sns"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
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
		log.Errorf(
			"Fraud detected for cart [%s], total confirmed by payment or manually = [%s], order total = [%s]",
			shoppingCart.ID, checkTotal.StringFixed(2), orderTotal.StringFixed(2),
		)
		return nil, fmt.Errorf("different totals: %s | %s", checkTotal.StringFixed(2), orderTotal.StringFixed(2))
	}

	order.InvoiceNumber = invoiceNumber
	order.OrderNumber = &orderNumber
	for i := 0; i < len(order.Items); i++ {
		// check if gift card, then create codes
		if order.Items[i].IsGiftCard {
			order.Items[i].GiftCardCode = lo.Times(order.Items[i].Amount, func(index int) string {
				coupon, err := db.NewGiftCard(ctx, order.Items[i].UnitPrice)
				if err != nil {
					log.Errorf("Failed to create coupon [%s]", err)
					// TODO: we don't expect this happens, so it's good for now
					return ""
				}
				log.Infof("Gift card code generated: %s = %s€", coupon.Code, coupon.Total)
				return coupon.Code
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

	sendSNSOrderCreated(ctx, order)

	return order, nil
}

func sendSNSOrderCreated(ctx context.Context, order *db.Order) {
	log := logger.FromContext(ctx)
	log.Infof("Send SNS message")

	topic, err := ssm.GetParameter(ctx, "/shop/sns/order/created/arn", false)
	if err != nil {
		log.Errorf("Failed to get topic arn: %s", err)
		return
	}
	c, err := mysns.NewClient(ctx)
	if err != nil {
		log.Errorf("Failed to init sns: %s", err)
		return
	}
	params := &sns.PublishInput{
		TopicArn:       aws.String(topic),
		Message:        aws.String(fmt.Sprintf("New order [%d]", order.OrderNumber)),
		MessageGroupId: aws.String(order.StoreKey),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"orderId": {
				DataType:    aws.String("String"),
				StringValue: aws.String(order.ID.Hex()),
			},
		},
	}
	if _, err := c.Publish(ctx, params); err != nil {
		log.Errorf("Failed to publish: %s", err)
		return
	}
	log.Infof("SNS message published")
}
