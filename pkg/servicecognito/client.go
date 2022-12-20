package servicecognito

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var (
	initClient sync.Once
	client     *cognitoidentityprovider.Client
)

func getClient(ctx context.Context) (*cognitoidentityprovider.Client, error) {
	var err error
	initClient.Do(func() {
		cfg, e := config.LoadDefaultConfig(ctx)
		if e != nil {
			err = e
			return
		}
		client = cognitoidentityprovider.NewFromConfig(cfg)
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
