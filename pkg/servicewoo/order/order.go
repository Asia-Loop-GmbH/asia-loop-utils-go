package order

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v7/pkg/servicewoo"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
)

func GetRefunds(ctx context.Context, id int) ([]servicewoo.Refund, error) {
	log := logger.FromContext(ctx)
	w, err := servicewoo.NewWoo(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	url := w.NewURL(ctx, fmt.Sprintf("/orders/%d/refunds", id))
	res, err := http.Get(url)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("failed to close response body: %s", err)
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		log.Errorf("status %d", res.StatusCode)
		return nil, fmt.Errorf("status %d", res.StatusCode)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	refunds := make([]servicewoo.Refund, 0)
	err = json.Unmarshal(resBody, &refunds)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	return refunds, nil
}

func Get(ctx context.Context, id int) (*servicewoo.Order, error) {
	log := logger.FromContext(ctx)
	w, err := servicewoo.NewWoo(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	url := w.NewURL(ctx, fmt.Sprintf("/orders/%d", id))
	res, err := http.Get(url)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("failed to close response body: %s", err)
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		log.Errorf("status %d", res.StatusCode)
		return nil, fmt.Errorf("status %d", res.StatusCode)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	wooOrder := new(servicewoo.Order)
	err = json.Unmarshal(resBody, wooOrder)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	return wooOrder, nil
}

func search(ctx context.Context, text string, page int, perPage int) ([]servicewoo.Order, error) {
	log := logger.FromContext(ctx)
	w, err := servicewoo.NewWoo(ctx)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	url := w.NewURL(ctx, fmt.Sprintf("/orders?search=%s&page=%d&per_page=%d", text, page, perPage))
	res, err := http.Get(url)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("failed to close response body: %s", err)
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		log.Errorf("status %d", res.StatusCode)
		return nil, fmt.Errorf("status %d", res.StatusCode)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	orders := make([]servicewoo.Order, 0)
	err = json.Unmarshal(resBody, &orders)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	return orders, nil
}

func SearchAll(ctx context.Context, text string, perPage int) ([]servicewoo.Order, error) {
	log := logger.FromContext(ctx)
	orders := make([]servicewoo.Order, 0)
	page := 1
	for true {
		log.Infof("Get orders [%s] page [%d]", text, page)
		os, err := search(ctx, text, page, perPage)
		if err != nil {
			log.Errorf("Failed to get orders: %s", err)
			return nil, errors.Wrap(err, "failed to get orders")
		}
		if len(os) == 0 {
			log.Infof("No result, stop")
			break
		}
		orders = append(orders, os...)
		page = page + 1
	}
	return orders, nil
}
