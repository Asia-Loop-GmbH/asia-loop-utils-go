package rest

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestGetSignInFromRequest(t *testing.T) {
	request := &events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{
			Identity: events.APIGatewayRequestIdentity{
				CognitoAuthenticationProvider: "cognito-idp.eu-central-1.amazonaws.com/eu-central-1_zubPX1a0g,cognito-idp.eu-central-1.amazonaws.com/eu-central-1_zubPX1a0g:CognitoSignIn:2671a429-b69b-453e-84b1-c6eb4d5cc551",
			},
		},
	}
	assert.Equal(t, "2671a429-b69b-453e-84b1-c6eb4d5cc551", GetSignInFromRequest(context.TODO(), request))
}

func TestGetSignInFromRequest_Malformed(t *testing.T) {
	request := &events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{
			Identity: events.APIGatewayRequestIdentity{
				CognitoAuthenticationProvider: "some invalid string",
			},
		},
	}
	assert.Equal(t, "", GetSignInFromRequest(context.TODO(), request))
}

func TestGetSignInFromRequest_Unsupported(t *testing.T) {
	request := &events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{
			Identity: events.APIGatewayRequestIdentity{
				CognitoAuthenticationProvider: "cognito-idp.eu-central-1.amazonaws.com/eu-central-1_zubPX1a0g,cognito-idp.eu-central-1.amazonaws.com/eu-central-1_zubPX1a0g:SomethingElse:2671a429-b69b-453e-84b1-c6eb4d5cc551",
			},
		},
	}
	assert.Equal(t, "", GetSignInFromRequest(context.TODO(), request))
}
