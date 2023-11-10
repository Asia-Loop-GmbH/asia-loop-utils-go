package auth

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/samber/lo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

var stagingIPs = []string{"35.205.49.4", "35.195.152.77", "34.77.114.185"}
var productionIPs = []string{"35.241.160.154", "35.241.180.107", "104.199.82.58", "34.79.19.218"}

var stageIPMap = map[string][]string{
	"dev":  stagingIPs,
	"pre":  productionIPs,
	"prod": productionIPs,
}

func ValidateCallerIP(ctx context.Context, request *events.APIGatewayProxyRequest) bool {
	log := logger.FromContext(ctx)
	callerIP := request.RequestContext.Identity.SourceIP
	log.Infof("Validate deliverect webhook call from [%s] to [%s]", callerIP, request.Path)

	stage := request.RequestContext.Stage
	allowedIPs, ok := stageIPMap[stage]
	if !ok {
		log.Errorf("Invalid environment [%s]", stage)
		return false
	}

	valid := lo.Contains(allowedIPs, callerIP)
	if !valid {
		log.Errorf("IP is not in this allow list: %v", allowedIPs)
		return false
	}

	return true
}
