package serviceadyen

import (
	"context"
	"fmt"
	"io"

	"github.com/adyen/adyen-go-api-library/v6/src/checkout"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

type RefundOptions struct {
	PSPRef      string
	Value       string
	MerchantRef string
}

func Refund(ctx context.Context, opts RefundOptions) (*checkout.PaymentRefundResource, error) {
	log := logger.FromContext(ctx)
	log.Infof("Cancel payment")

	client, err := newClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init adyen client")
	}

	amount := decimal.RequireFromString(opts.Value)
	amountInt := amount.Mul(decimal.NewFromInt(100)).IntPart()
	res, httpRes, err := client.Checkout.PaymentsPaymentPspReferenceRefunds(opts.PSPRef, &checkout.CreatePaymentRefundRequest{
		Amount: checkout.Amount{
			Currency: currencyEUR,
			Value:    amountInt,
		},
		MerchantAccount: accountECOM,
		Reference:       opts.MerchantRef,
	}, ctx)

	if err != nil {
		log.Errorf("Failed to refund payment [%+v]: %s", opts, err)
		return nil, errors.Wrap(err, "failed to refund payment")
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

	log.Infof("Refund succeeded: %+v", opts)
	return &res, nil
}
