package rest

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v4/pkg/shop/db"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/rest"
)

func RequireAdmin(ctx context.Context, request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	log := logger.FromContext(ctx)
	log.Infof("This action requires admin role")

	colRoles, err := db.CollectionUsers(ctx)
	if err != nil {
		return rest.ResponseError(ctx, http.StatusInternalServerError, request, errors.Wrap(err, "failed to init db collection"))
	}

	username := GetSignInFromRequest(ctx, request)
	find := colRoles.FindOne(ctx, bson.M{"user": username})
	user := new(db.User)
	err = find.Decode(user)
	if err == mongo.ErrNoDocuments {
		log.Infof("%+v", request.RequestContext.Identity)
		return rest.ResponseError(ctx, http.StatusForbidden, request, errors.Wrapf(err, "user [%s] not found", username))
	}
	if err != nil {
		return rest.ResponseError(ctx, http.StatusInternalServerError, request, errors.Wrapf(err, "failed to find user [%s]", username))
	}

	if !lo.ContainsBy(user.Roles, func(item db.RoleEntry) bool {
		return item.Name == db.RoleAdmin
	}) {
		log.Infof("%+v", request.RequestContext.Identity)
		log.Errorf("User [%s] doesn't have admin role", username)
		return rest.Response(ctx, http.StatusForbidden, request, nil)
	}

	return nil
}

func GetSignInFromRequest(ctx context.Context, request *events.APIGatewayProxyRequest) string {
	log := logger.FromContext(ctx)
	log.Infof("Get signin from request")

	auth := request.RequestContext.Identity.CognitoAuthenticationProvider
	parts := strings.Split(auth, ":")
	if len(parts) < 3 {
		log.Warnf("Malformed auth: %s", auth)
		return ""
	}
	if parts[len(parts)-2] != "CognitoSignIn" {
		log.Warnf("Unsupported auth: %s", auth)
		return ""
	}
	signIn := parts[len(parts)-1]
	log.Infof("Sign in user: %s", signIn)
	return signIn
}
