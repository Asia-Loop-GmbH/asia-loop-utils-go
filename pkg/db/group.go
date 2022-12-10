package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/pkg/mongodb"
)

const colGroups = "groups"

func CollectionGroups(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/mongo/database", false)
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colGroups)
}

var (
	GroupObjectIDs = []string{
		"orders", "driver", "store",
	}
)

type Group struct {
	ID         primitive.ObjectID   `bson:"_id" json:"id"`
	Orders     []primitive.ObjectID `bson:"orders" json:"orders"`
	RouteOrder []int                `bson:"routeOrder" json:"routeOrder"`
	Number     string               `bson:"number" json:"number"`
	Finalized  bool                 `bson:"finalized" json:"finalized"`
	Delivered  bool                 `bson:"delivered" json:"delivered"`
	Driver     primitive.ObjectID   `bson:"driver" json:"driver"`
	DriverName string               `bson:"driverName" json:"driverName"`
	Store      primitive.ObjectID   `bson:"store" json:"store"`
	CreatedAt  time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time            `bson:"updatedAt" json:"updatedAt"`
}
