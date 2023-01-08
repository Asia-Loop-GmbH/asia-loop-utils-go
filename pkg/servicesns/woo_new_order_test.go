package servicesns_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/pkg/servicesns"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v2/pkg/context"
)

func TestPublishWooNewOrder(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	err := servicesns.PublishWooNewOrder(ctx, &servicesns.EventWooNewOrderData{
		ID: 1234,
	})
	assert.NoError(t, err)
}
