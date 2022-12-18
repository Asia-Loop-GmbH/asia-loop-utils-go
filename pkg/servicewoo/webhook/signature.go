package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-lambda-go/events"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
)

const (
	signatureHeader = "X-WC-Webhook-Signature"
)

// ValidateRequest
// returns isPing, isValid, err
func ValidateRequest(ctx context.Context, request *events.APIGatewayProxyRequest, webhookID int) (bool, bool, error) {
	log := logger.FromContext(ctx)
	pingBody := fmt.Sprintf("webhook_id=%d", webhookID)

	if request.Body == pingBody {
		log.Infof("body = [%s], this is a ping request from woocommerce, skip processing", request.Body)
		// ping=true, valid=true, err=nil
		return true, true, nil
	}

	signature := request.Headers[signatureHeader]
	log.Infof(request.Body)
	log.Infof(fmt.Sprintf("%v", request.Headers))

	secret, err := ssm.GetParameter(ctx, "/shop/woo/webhook/secret", true)
	if err != nil {
		log.Errorf("failed to get webhook secret from ssm: %s", err)
		// ping=false, valid=false, err=err
		return false, false, err
	}

	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(request.Body))
	sha := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	if sha != signature {
		log.Errorf("wrong signature: expected = %s, got = %s", sha, signature)
		// ping=false, valid=false, err=nil
		return false, false, nil
	}

	// ping=false, valid=true, err=nil
	return false, true, nil
}
