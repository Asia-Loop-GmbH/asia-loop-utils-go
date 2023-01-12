package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/mongodb"
)

const (
	RoleAdmin         = "Admin"
	RoleStoreAdmin    = "Store Admin"
	RoleStoreStandard = "Store Standard"
)

type RoleEntry struct {
	Name  string  `bson:"name" json:"name"`
	Store *string `bson:"store,omitempty" json:"store,omitempty"`
}

type User struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	User  string             `bson:"user" json:"user"`
	Email string             `bson:"email" json:"email"`
	Roles []RoleEntry        `bson:"roles" json:"roles"`
}

const colUsers = "users"

func CollectionUsers(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/shop/mongo/database", false)
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colUsers)
}
