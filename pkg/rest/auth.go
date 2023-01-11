package rest

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v2/pkg/shop/db"
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

	username := request.RequestContext.Identity.User
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
