package serviceadyen

import (
	"context"

	"github.com/adyen/adyen-go-api-library/v6/src/adyen"
	"github.com/adyen/adyen-go-api-library/v6/src/common"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
)

func newClient(ctx context.Context) (*adyen.APIClient, error) {
	log := logger.FromContext(ctx)

	log.Infof("new adyen client")
	apiEnv, err := ssm.GetParameter(ctx, "/adyen/env", false)
	if err != nil {
		return nil, err
	}
	apiKey, err := ssm.GetParameter(ctx, "/adyen/key", true)
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
