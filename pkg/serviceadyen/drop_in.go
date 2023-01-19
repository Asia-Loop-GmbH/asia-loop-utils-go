package serviceadyen

import (
	"context"
	"fmt"
	"io"

	"github.com/adyen/adyen-go-api-library/v6/src/checkout"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v3/pkg/shop/cart"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/random"
)

func NewDropInPayment(ctx context.Context, shoppingCart *cart.PublicCart, returnURL string) (*checkout.CreateCheckoutSessionResponse, error) {
	log := logger.FromContext(ctx)
	log.Infof("new drop in payment for cart [%s]", shoppingCart.ID.Hex())
	client, err := newClient(ctx)
	if err != nil {
		return nil, err
	}

	amount := decimal.RequireFromString(shoppingCart.Summary.Total.Value)
	amountInt := amount.Mul(decimal.NewFromInt(100)).IntPart()

	res, httpRes, err := client.Checkout.Sessions(&checkout.CreateCheckoutSessionRequest{
		Amount:                 checkout.Amount{Currency: currencyEUR, Value: amountInt},
		CountryCode:            countryDE,
		MerchantAccount:        accountECOM,
		MerchantOrderReference: shoppingCart.ID.Hex(),
		Reference:              random.String(10, lo.AlphanumericCharset),
		ReturnUrl:              returnURL,
		// TODO: should we send more data to adyen
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
