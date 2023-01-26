package globalconfig

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v7/pkg/db"
)

func Get(ctx context.Context) (*db.GlobalConfig, error) {
	colGlbCfg, err := db.CollectionGlobalConfig(ctx)
	if err != nil {
		return nil, err
	}
	find := colGlbCfg.FindOne(context.Background(), bson.M{})
	globalConfig := new(db.GlobalConfig)
	err = find.Decode(globalConfig)
	if err != nil {
		return nil, err
	}
	return globalConfig, nil
}
