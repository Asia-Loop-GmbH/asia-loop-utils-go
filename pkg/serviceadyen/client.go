package serviceadyen

import (
	"context"
	"fmt"

	"github.com/adyen/adyen-go-api-library/v5/src/adyen"
	"github.com/adyen/adyen-go-api-library/v5/src/common"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/ssm"
	mycontext "github.com/nam-truong-le/lambda-utils-go/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
)

var (
	envMap = map[string]common.Environment{
		"dev":  common.TestEnv,
		"pre":  common.LiveEnv,
		"prod": common.LiveEnv,
	}
)

func newClient(ctx context.Context) (*adyen.APIClient, error) {
	log := logger.FromContext(ctx)
	stage, ok := ctx.Value(mycontext.FieldStage).(string)
	if !ok {
		return nil, fmt.Errorf("undefined stage in context")
	}

	log.Infof("new adyen client: %s", stage)
	apiKey, err := ssm.GetParameter(ctx, "/adyen/key", true)
	if err != nil {
		return nil, err
	}
	environment, ok := envMap[stage]
	if !ok {
		return nil, fmt.Errorf("no adyen environment config found for stage: %s", stage)
	}
	if stage == "dev" {
		return adyen.NewClient(&common.Config{
			ApiKey:      apiKey,
			Environment: environment,
		}), nil
	}

	return adyen.NewClient(&common.Config{
		ApiKey:                apiKey,
		Environment:           environment,
		LiveEndpointURLPrefix: livePrefix,
	}), nil
}
