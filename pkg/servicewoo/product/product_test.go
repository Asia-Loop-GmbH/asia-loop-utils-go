package product_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v9/pkg/servicewoo/product"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
)

func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.TODO(), mycontext.FieldStage, "dev")
	products, err := product.Get(ctx)
	assert.NoError(t, err)
	assert.True(t, len(products) > 0)
}
