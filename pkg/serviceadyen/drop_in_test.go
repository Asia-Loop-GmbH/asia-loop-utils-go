package serviceadyen_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/pkg/random"
	"github.com/asia-loop-gmbh/asia-loop-utils-go/pkg/serviceadyen"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v2/pkg/context"
)

func TestNewDropInPayment_Success(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	response, err := serviceadyen.NewDropInPayment(
		ctx,
		"10.23",
		random.String(10, true, true, true),
		"https://admin2-dev.asia-loop.com",
	)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}
