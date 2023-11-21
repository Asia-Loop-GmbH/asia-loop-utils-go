package serviceadyen

import (
	"context"

	"github.com/adyen/adyen-go-api-library/v8/src/adyen"
	"github.com/adyen/adyen-go-api-library/v8/src/common"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

func newClient(ctx context.Context) (*adyen.APIClient, error) {
	log := logger.FromContext(ctx)

	log.Infof("new adyen client")
	apiEnv, err := secretsmanager.GetParameter(ctx, "/adyen/env")
	if err != nil {
		return nil, err
	}
	apiKey, err := secretsmanager.GetParameter(ctx, "/adyen/key")
	if err != nil {
		return nil, err
	}

	stage := ctx.Value(mycontext.FieldStage).(string)
	if stage == "dev" {
		return adyen.NewClient(&common.Config{
			ApiKey:      apiKey,
			Environment: common.Environment(apiEnv),
		}), nil
	}

	return adyen.NewClient(&common.Config{
		ApiKey:                apiKey,
		Environment:           common.Environment(apiEnv),
		LiveEndpointURLPrefix: livePrefix,
	}), nil
}
