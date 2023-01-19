package servicesns_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v3/pkg/servicesns"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
)

func TestPublishOrderCreated(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	err := servicesns.PublishOrderCreated(ctx, &servicesns.EventOrderCreatedData{
		OrderID: "POS-810052",
	})
	assert.NoError(t, err)
}
