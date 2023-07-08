package pdf

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

type Request struct {
	Content string `json:"content"`
}

type Response struct {
	Base64 string `json:"base64"`
}

func NewFromHTML(ctx context.Context, html string) ([]byte, error) {
	log := logger.FromContext(ctx)
	toolsURL, err := secretsmanager.GetParameter(ctx, "/external/tools/url")
	if err != nil {
		return nil, err
	}
	apiKey, err := secretsmanager.GetParameter(ctx, "/external/tools/apikey")
	if err != nil {
		return nil, err
	}

	body := Request{Content: html}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	client := http.Client{Timeout: time.Second * 30}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/pdf/v1/generate", toolsURL),
		bytes.NewReader(bodyJSON),
	)
	req.Header.Add("X-Api-Key", apiKey)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("Failed to close response body: %s", err)
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		log.Errorf("HTTP code: %d", res.StatusCode)
		return nil, fmt.Errorf("failed to convert html to pdf, http code: %d", res.StatusCode)
	}

	resBodyRaw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	resBody := new(Response)
	err = json.Unmarshal(resBodyRaw, resBody)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(resBody.Base64)
}
