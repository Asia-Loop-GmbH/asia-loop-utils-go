package db

import (
	"context"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/mongodb"
)

const colStores = "stores"

func CollectionStores(ctx context.Context) (*mongo.Collection, error) {
	database, err := ssm.GetParameter(ctx, "/shop/mongo/database", false)
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colStores)
}

type Store struct {
	ID                   primitive.ObjectID `bson:"_id" json:"id"`
	Email                string             `bson:"email" json:"email"`
	Telephone            string             `bson:"telephone" json:"telephone"`
	Name                 string             `bson:"name" json:"name"`
	Key                  string             `bson:"key" json:"key"`
	Address              string             `bson:"address" json:"address"`
	Owner                string             `bson:"owner" json:"owner"`
	BusinessRegistration string             `bson:"businessRegistration" json:"businessRegistration"`
	TaxNumber            string             `bson:"taxNumber" json:"taxNumber"`
	MBW                  map[string]string  `bson:"mbw" json:"mbw"`
	Slots                []map[string]bool  `bson:"slots" json:"slots"` // week starts with Sunday = index 0
	Holidays             []string           `bson:"holidays" json:"holidays"`
	CreatedAt            time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt            time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func (s *Store) GetSlots(date string) []string {
	if lo.Contains(s.Holidays, date) {
		return make([]string, 0)
	}

	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return make([]string, 0)
	}

	weekDay := int(d.Weekday())
	slots := s.Slots[weekDay]
	return lo.Filter(lo.Keys(slots), func(item string, index int) bool {
		return slots[item]
	})
}
