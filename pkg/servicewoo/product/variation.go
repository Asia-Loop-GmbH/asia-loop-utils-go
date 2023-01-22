package product

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v6/pkg/servicewoo"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
)

func GetVariation(ctx context.Context, productID int) ([]servicewoo.ProductVariation, error) {
	log := logger.FromContext(ctx)
	woo, err := servicewoo.NewWoo(ctx)
	if err != nil {
		return nil, err
	}
	url := woo.NewURL(ctx, fmt.Sprintf("/products/%d/variations", productID))
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("Failed to close body: %s", err)
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	variations := make([]servicewoo.ProductVariation, 0)
	err = json.Unmarshal(body, &variations)
	if err != nil {
		return nil, err
	}
	return variations, nil
}
