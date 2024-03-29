package db

import (
	"context"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mongodb"
)

const colStores = "stores"

func CollectionStores(ctx context.Context) (*mongo.Collection, error) {
	database, err := secretsmanager.GetParameter(ctx, "/shop/mongo/database")
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colStores)
}

type DeliverectIntegration struct {
	AccountID  string `bson:"accountId" json:"accountId"`
	LocationID string `bson:"locationId" json:"locationId"`
}

type Store struct {
	ID                    primitive.ObjectID         `bson:"_id" json:"id"`
	Email                 string                     `bson:"email" json:"email"`
	Telephone             string                     `bson:"telephone" json:"telephone"`
	Name                  string                     `bson:"name" json:"name"`
	Key                   string                     `bson:"key" json:"key"`
	Address               string                     `bson:"address" json:"address"`
	Owner                 string                     `bson:"owner" json:"owner"`
	BusinessRegistration  string                     `bson:"businessRegistration" json:"businessRegistration"`
	TaxNumber             string                     `bson:"taxNumber" json:"taxNumber"`
	MBW                   map[string]string          `bson:"mbw" json:"mbw"`
	MBWAllowOnlyCities    map[string][]string        `bson:"mbwAllowOnlyCities" json:"mbwAllowOnlyCities,omitempty"`
	Slots                 []map[string]bool          `bson:"slots" json:"slots"` // week starts with Sunday = index 0
	Holidays              []string                   `bson:"holidays" json:"holidays"`
	SpecialDays           map[string]map[string]bool `bson:"specialDays,omitempty" json:"specialDays,omitempty"`
	CreatedAt             time.Time                  `bson:"createdAt" json:"createdAt"`
	UpdatedAt             time.Time                  `bson:"updatedAt" json:"updatedAt"`
	DeliverectIntegration *DeliverectIntegration     `bson:"deliverectIntegration,omitempty" json:"deliverectIntegration,omitempty"`
}

func (s *Store) GetSlots(date string) []string {
	if lo.Contains(s.Holidays, date) {
		return make([]string, 0)
	}

	if specialDay, ok := s.SpecialDays[date]; ok {
		return lo.Filter(lo.Keys(specialDay), func(item string, index int) bool {
			return specialDay[item]
		})
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
