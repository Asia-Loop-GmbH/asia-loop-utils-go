package product

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/pkg/servicewoo"
	"github.com/nam-truong-le/lambda-utils-go/v2/pkg/logger"
)

const (
	perPage = 100
)

func Get(ctx context.Context) ([]servicewoo.Product, error) {
	log := logger.FromContext(ctx)
	log.Infof("get all products from woo")
	page := 1
	result := make([]servicewoo.Product, 0)
	for true {
		ps, err := getPage(ctx, page)
		if err != nil {
			return nil, err
		}
		if len(ps) == 0 {
			break
		}
		for _, p := range ps {
			if p.Status == "publish" {
				result = append(result, p)
			}
		}
		page++
	}
	return result, nil
}

func getPage(ctx context.Context, page int) ([]servicewoo.Product, error) {
	log := logger.FromContext(ctx)
	log.Infof("get products from woo page [%d], per_page [%d]", page, perPage)
	woo, err := servicewoo.NewWoo(ctx)
	if err != nil {
		return nil, err
	}

	url := woo.NewURL(ctx, fmt.Sprintf("/products?page=%d&per_page=%d", page, perPage))
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("failed to close response body: %s", err)
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("woo returns [%s]", res.Status)
	}
	products := make([]servicewoo.Product, 0)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &products)
	return products, nil
}
