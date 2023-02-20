package serviceadyen

import (
	"context"
	"fmt"
	"io"

	"github.com/adyen/adyen-go-api-library/v6/src/checkout"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v8/pkg/shop/db"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/random"
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

	amount := decimal.RequireFromString(order.Summary.Total.Value)
	amountInt := amount.Mul(decimal.NewFromInt(100)).IntPart()

	res, httpRes, err := client.Checkout.Sessions(&checkout.CreateCheckoutSessionRequest{
		Amount:                 checkout.Amount{Currency: currencyEUR, Value: amountInt},
		CountryCode:            countryDE,
		MerchantAccount:        accountECOM,
		MerchantOrderReference: order.ID.Hex(),
		Reference:              random.String(10, lo.AlphanumericCharset),
		ReturnUrl:              returnURL,
		// TODO: should we send more data to adyen
		LineItems: lo.ToPtr(lo.Map(order.Items, func(item db.OrderItem, idx int) checkout.LineItem {
			return checkout.LineItem{
				Id:                 fmt.Sprintf("Item #%d", idx),
				AmountExcludingTax: adyenUnitPrice(item.Net, item.Amount),
				AmountIncludingTax: adyenUnitPrice(item.Total, item.Amount),
				Description:        item.Name,
				Quantity:           int64(item.Amount),
				TaxAmount:          adyenUnitPrice(item.Tax, item.Amount),
				TaxPercentage:      adyenTaxPercentage(item.TaxClass),
			}
		})),
		ShopperEmail:     cartCheckout.Email,
		ShopperReference: random.String(10, lo.UpperCaseLettersCharset), // TODO: is there anyway to have a better shopper reference?
	})

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

	res, httpRes, err := client.Checkout.PaymentsDetails(&checkout.DetailsRequest{
		Details: checkout.PaymentCompletionDetails{
			RedirectResult: redirectResult,
		},
	})
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
