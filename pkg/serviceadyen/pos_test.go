package serviceadyen_test

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v7/pkg/serviceadyen"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/random"
)

func TestNewTender(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	orderId := random.String(6, lo.AlphanumericCharset)
	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	err := serviceadyen.NewTender(ctx, "S1F2-000158213300585", orderId, 10.12)
	assert.Nil(t, err)
}
