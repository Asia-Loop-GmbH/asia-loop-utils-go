package serviceadyen_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v8/pkg/serviceadyen"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
)

func TestRefund(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.Background(), mycontext.FieldStage, "dev")
	res, err := serviceadyen.Refund(ctx, serviceadyen.RefundOptions{
		RefundRef: "",
		PSPRef:    "G8WVVDD5HV5X8N82",
		Value:     "16.45",
		Items:     nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1645), res.Amount.Value)
}
