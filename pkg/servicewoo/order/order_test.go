package order_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/pkg/servicewoo/order"
	mycontext "github.com/nam-truong-le/lambda-utils-go/pkg/context"
)

func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	o, err := order.Get(ctx, 123)
	assert.NoError(t, err)
	assert.Equal(t, 123, o.ID)
}

func TestGetRefunds(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	refunds, err := order.GetRefunds(ctx, 245)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(refunds))
	assert.Equal(t, 262, refunds[0].ID)
}
