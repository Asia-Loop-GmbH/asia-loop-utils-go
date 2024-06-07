package serviceadyen

import (
	"context"
	"fmt"
	"io"

	"github.com/adyen/adyen-go-api-library/v8/src/checkout"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v9/pkg/shop/db"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/random"
)

func adyenUnitPrice(sum string, amount int) int64 {
	return decimal.RequireFromString(sum).Div(decimal.NewFromInt(int64(amount))).Mul(decimal.NewFromInt(100)).IntPart()
}

func adyenTaxPercentage(taxClass string) int64 {
	switch taxClass {
	case db.TaxClassStandard:
		return 1900
	case db.TaxClassZero:
		return 0
	default:
		return 700
	}
}

// NewDropInPayment contains order that is converted from shopping cart and hence doesn't have checkout data. We must
// pass checkout data separately.
func NewDropInPayment(ctx context.Context, order *db.Order, cartCheckout *db.CartCheckout, returnURL string) (*checkout.CreateCheckoutSessionResponse, error) {
	log := logger.FromContext(ctx)
	log.Infof("new drop in payment for cart [%s]", order.ID.Hex())
	client, err := newClient(ctx)
	if err != nil {
		return nil, err
	}

	checkoutService := client.Checkout()
	amount := decimal.RequireFromString(order.Summary.Total.Value)
	amountInt := amount.Mul(decimal.NewFromInt(100)).IntPart()

	createCheckoutSession := checkout.CreateCheckoutSessionRequest{
		Amount:                 checkout.Amount{Currency: currencyEUR, Value: amountInt},
		CountryCode:            lo.ToPtr(countryDE),
		MerchantAccount:        accountECOM,
		MerchantOrderReference: lo.ToPtr(order.ID.Hex()),
		Reference:              random.String(10, lo.AlphanumericCharset),
		ReturnUrl:              returnURL,
		LineItems: lo.Map(order.Items, func(item db.OrderItem, idx int) checkout.LineItem {
			return checkout.LineItem{
				Id:                 lo.ToPtr(fmt.Sprintf("Item #%d", idx)),
				AmountExcludingTax: lo.ToPtr(adyenUnitPrice(item.Net, item.Amount)),
				AmountIncludingTax: lo.ToPtr(adyenUnitPrice(item.Total, item.Amount)),
				Description:        lo.ToPtr(item.Name),
				Quantity:           lo.ToPtr(int64(item.Amount)),
				TaxAmount:          lo.ToPtr(adyenUnitPrice(item.Tax, item.Amount)),
				TaxPercentage:      lo.ToPtr(adyenTaxPercentage(item.TaxClass)),
			}
		}),
		ShopperEmail: lo.ToPtr(cartCheckout.Email),
	}
	idempotencyKey := uuid.New().String()
	res, httpRes, err := checkoutService.PaymentsApi.Sessions(ctx,
		checkoutService.PaymentsApi.SessionsInput().IdempotencyKey(idempotencyKey).CreateCheckoutSessionRequest(createCheckoutSession),
	)

	if err != nil {
		log.Errorf("Failed to create payment session: %s", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("Failed to close http response body: %s", err)
		}
	}(httpRes.Body)

	if httpRes.StatusCode >= 300 {
		responseBody, err := io.ReadAll(httpRes.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("adyen request error %s: %s", httpRes.Status, string(responseBody))
	}

	log.Infof("%+v", res)
	return &res, nil
}

func ProcessRedirect(ctx context.Context, sessionID, redirectResult string) (*checkout.PaymentDetailsResponse, error) {
	log := logger.FromContext(ctx)
	log.Infof("Update details for session [%s]", sessionID)
	client, err := newClient(ctx)
	if err != nil {
		return nil, err
	}
	checkoutService := client.Checkout()

	res, httpRes, err := checkoutService.PaymentsApi.PaymentsDetails(ctx, checkoutService.PaymentsApi.PaymentsDetailsInput().IdempotencyKey(sessionID).PaymentDetailsRequest(checkout.PaymentDetailsRequest{
		Details: checkout.PaymentCompletionDetails{
			RedirectResult: lo.ToPtr(redirectResult),
		},
	}))
	if err != nil {
		log.Errorf("Failed to update payment details: %s", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("Failed to close http response body: %s", err)
		}
	}(httpRes.Body)

	if httpRes.StatusCode >= 300 {
		responseBody, err := io.ReadAll(httpRes.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("dyen request error %s: %s", httpRes.Status, string(responseBody))
	}

	log.Infof("%+v", res)
	return &res, nil
}
