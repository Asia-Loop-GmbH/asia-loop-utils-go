package servicesns_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/pkg/servicesns"
	mycontext "github.com/nam-truong-le/lambda-utils-go/pkg/context"
)

func TestPublishOrderPOSPaid(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.Background(), mycontext.FieldStage, "dev")
	err := servicesns.PublishOrderPOSPaid(ctx, &servicesns.EventOrderPOSPaidData{
		OrderID: "POS-810052",
	})
	assert.NoError(t, err)
}