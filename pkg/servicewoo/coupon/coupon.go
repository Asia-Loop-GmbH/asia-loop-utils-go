package coupon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v4/pkg/servicewoo"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
)

func IsValidAndHasEnough(ctx context.Context, code, appliedAmount string) bool {
	log := logger.FromContext(ctx)
	log.Infof("check coupon valid and has enough amount: %s", code)
	coupon, err := GetCouponByCode(ctx, code)
	if err != nil {
		log.Errorf("could not get coupon '%s': %s", code, err)
		return false
	}
	current, err := decimal.NewFromString(coupon.Amount)
	if err != nil {
		return false
	}
	toUse, err := decimal.NewFromString(appliedAmount)
	if err != nil {
		return false
	}
	return current.Cmp(toUse) > 0
}

func GetCouponByCode(ctx context.Context, code string) (*servicewoo.Coupon, error) {
	log := logger.FromContext(ctx)
	log.Infof("get coupon: %s", code)
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, fmt.Errorf("blank coupon code")
	}
	serviceWoo, err := servicewoo.NewWoo(ctx)
	if err != nil {
		return nil, err
	}

	response, err := http.Get(serviceWoo.NewURL(ctx, fmt.Sprintf("/coupons?code=%s", code)))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("failed to close response body: %s", err)
		}
	}(response.Body)
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	coupons := make([]servicewoo.Coupon, 0)
	if err := json.Unmarshal(responseBody, &coupons); err != nil {
		return nil, err
	}

	if len(coupons) == 0 {
		return nil, fmt.Errorf("coupon code '%s' not found", code)
	}
	if len(coupons) > 1 {
		return nil, fmt.Errorf("multiple coupon codes found for '%s'", code)
	}
	return &coupons[0], nil
}

func UpdateCouponByCode(ctx context.Context, code, amount string) error {
	log := logger.FromContext(ctx)
	log.Infof("update coupon %s: %s", code, amount)
	coupon, err := GetCouponByCode(ctx, code)
	if err != nil {
		return err
	}
	toUse, err := decimal.NewFromString(amount)
	if err != nil {
		return err
	}
	currentAmount, err := decimal.NewFromString(coupon.Amount)
	if err != nil {
		return err
	}
	newAmount := currentAmount.Sub(toUse)

	updateCoupon := servicewoo.Coupon{
		Amount: newAmount.StringFixed(2),
	}

	serviceWoo, err := servicewoo.NewWoo(ctx)
	if err != nil {
		return err
	}

	requestBody, err := json.Marshal(updateCoupon)
	if err != nil {
		return err
	}
	log.Infof("update coupon request body: %s", string(requestBody))

	url := serviceWoo.NewURL(ctx, fmt.Sprintf("/coupons/%d", coupon.ID))
	log.Infof("PUT -> %s", url)
	httpClient := http.Client{}
	request, err := http.NewRequest(
		http.MethodPut,
		url,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("failed to close response body: %s", err)
		}
	}(response.Body)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	log.Infof("response: %s", string(responseBody))

	if response.StatusCode >= 300 {
		return fmt.Errorf("could not update coupon, error '%d': %s", response.StatusCode, string(responseBody))
	}

	return nil
}
