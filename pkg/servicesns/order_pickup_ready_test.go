package servicesns_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/pkg/servicesns"
	mycontext "github.com/nam-truong-le/lambda-utils-go/pkg/context"
)

func TestPublishOrderPickupReady(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.Background(), mycontext.FieldStage, "dev")
	err := servicesns.PublishOrderPickupReady(ctx, &servicesns.EventOrderPickupReadyData{
		OrderID: "POS-810052",
		InTime:  "10 Minuten",
	})
	assert.NoError(t, err)
}
