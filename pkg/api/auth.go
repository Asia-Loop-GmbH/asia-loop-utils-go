package api

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
)

const (
	HeaderAPIKey = "x-al-api-key"
)

func AuthorizeAnalyticsRequest(ctx context.Context, request *events.APIGatewayProxyRequest) error {
	log := logger.FromContext(ctx)

	log.Infof("authorize analytics request [%s] [%s]", request.HTTPMethod, request.Path)

	keyInReq, hasKey := request.Headers[HeaderAPIKey]
	if !hasKey {
		return errors.New("missing API key")
	}

	key, err := ssm.GetParameter(ctx, "/analytics/key", true)
	if err != nil {
		return err
	}

	if keyInReq != key {
		return errors.New("unauthorized")
	}

	return nil
}
