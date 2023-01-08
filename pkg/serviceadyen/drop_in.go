package serviceadyen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/adyen/adyen-go-api-library/v6/src/checkout"
	"github.com/shopspring/decimal"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/pkg/api"
	"github.com/nam-truong-le/lambda-utils-go/v2/pkg/logger"
)

type SessionResponse struct {
	ID          string `json:"id"`
	SessionData string `json:"sessionData"`
}

func NewDropInPayment(ctx context.Context, value, ref, returnURL string) (*api.PaymentDropInResponse, error) {
	log := logger.FromContext(ctx)
	log.Infof("new drop in payment for order [%s]", ref)
	client, err := newClient(ctx)
	if err != nil {
		return nil, err
	}

	amount, err := decimal.NewFromString(value) // TODO: only valid for corporate customers!!!
	if err != nil {
		return nil, err
	}
	amountInt := amount.Mul(decimal.NewFromInt(100)).IntPart()

	req := &checkout.PaymentSetupRequest{
		Amount:          checkout.Amount{Currency: currencyEUR, Value: amountInt},
		MerchantAccount: accountECOM,
		ReturnUrl:       returnURL,
		Reference:       ref,
		CountryCode:     countryDE,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/v68/sessions", client.GetConfig().CheckoutEndpoint)
	log.Printf("POST -> %s", url)
	httpClient := http.Client{}
	postRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	postRequest.Header.Set("x-API-key", client.GetConfig().ApiKey)
	response, err := httpClient.Do(postRequest)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("failed to close response body: %s", err)
		}
	}(response.Body)

	if response.StatusCode >= 300 {
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("adyen request error %s: %s", response.Status, string(responseBody))
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	newSession := new(SessionResponse)
	if err := json.Unmarshal(responseBody, newSession); err != nil {
		return nil, err
	}
	result := api.PaymentDropInResponse{
		ID:          newSession.ID,
		SessionData: newSession.SessionData,
	}
	log.Infof("%v", result)
	return &result, nil
}
