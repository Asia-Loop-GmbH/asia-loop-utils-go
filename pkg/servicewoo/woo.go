package servicewoo

import (
	"context"
	"fmt"
	"strings"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

type Woo struct {
	URL    string
	Key    string
	Secret string
}

func NewWoo(ctx context.Context) (*Woo, error) {
	log := logger.FromContext(ctx)
	log.Infof("read woo information")
	shopUrl, err := secretsmanager.GetParameter(ctx, "/shop/woo/url")
	if err != nil {
		return nil, err
	}
	wooKey, err := secretsmanager.GetParameter(ctx, "/shop/woo/key")
	if err != nil {
		return nil, err
	}
	wooSecret, err := secretsmanager.GetParameter(ctx, "/shop/woo/secret")
	if err != nil {
		return nil, err
	}

	return &Woo{shopUrl, wooKey, wooSecret}, nil
}

func (w *Woo) NewURL(ctx context.Context, url string) string {
	return w.newURL(ctx, url, "/wp-json/wc/v3")
}

func (w *Woo) NewURLAsiaLoop(ctx context.Context, url string) string {
	return w.newURL(ctx, url, "/wp-json/asialoop-api")
}

func (w *Woo) newURL(ctx context.Context, url string, api string) string {
	log := logger.FromContext(ctx)
	log.Infof("prepare woo url: %s", url)
	connector := "?"
	if strings.Contains(url, "?") {
		connector = "&"
	}
	result := fmt.Sprintf(
		"%s%s%s%sconsumer_key=%s&consumer_secret=%s",
		w.URL, api, url, connector, w.Key, w.Secret,
	)
	log.Infof("final woo url: %s", result)
	return result
}
