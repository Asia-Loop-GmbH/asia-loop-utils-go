package product_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v6/pkg/servicewoo/product"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
)

func TestGetVariation(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	variations, err := product.GetVariation(ctx, 24)
	assert.NoError(t, err)
	assert.True(t, len(variations) > 0)
}
