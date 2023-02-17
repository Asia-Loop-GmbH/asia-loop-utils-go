package servicesns_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v8/pkg/servicesns"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
)

func TestPublishOrderFinalized(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	id, _ := primitive.ObjectIDFromHex("622a5275b73e4d6262fd8acf")
	err := servicesns.PublishOrderFinalized(ctx, &servicesns.EventOrderFinalizedData{
		ID: id,
	})
	assert.NoError(t, err)
}
