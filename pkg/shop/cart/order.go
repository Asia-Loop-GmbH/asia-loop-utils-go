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

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v9/pkg/orderutils"
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v9/pkg/shop/coupon"
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v9/pkg/shop/db"
	mysns "github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/sns"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/random"
)

func CreateOrder(ctx context.Context, shoppingCart *db.Cart, confirmedTotal string) (*db.Order, error) {
	log := logger.FromContext(ctx)
	log.Infof("Cart was paid, will create final order")

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
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now
	for i := 0; i < len(order.Items); i++ {
		// check if gift card, then create codes
		if order.Items[i].IsGiftCard {
			order.Items[i].GiftCardCode = lo.Times(order.Items[i].Amount, func(index int) string {
				newCoupon, err := db.NewMehrzweckCoupon(ctx, order.Items[i].UnitPrice)
				if err != nil {
					log.Errorf("Failed to create coupon [%s]", err)
					// TODO: we don't expect this happens, so it's good for now
					return ""
				}
				log.Infof("Gift card code generated: %s = %s€", newCoupon.Code, newCoupon.Total)
				return newCoupon.Code
			})
		}

		// check if coupon, then update coupon value
		if order.Items[i].SKU == db.CouponSKU {
			err := coupon.UpdateCouponByOrderItem(ctx, order.ID.Hex(), order.Items[i])
			if err != nil {
				log.Errorf("Failed to update coupon %s: %s", order.Items[i].Name, err)
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
