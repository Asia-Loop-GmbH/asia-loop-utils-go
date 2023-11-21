package db_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v9/pkg/db"
	commoncontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mongodb"
)

func TestNextByStage(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.Background(), commoncontext.FieldStage, "dev")
	defer mongodb.Disconnect(ctx)
	now, err := db.Next(ctx, "test")
	assert.NoError(t, err)
	assert.True(t, now > 0)

	next, err := db.Next(ctx, "test")
	assert.NoError(t, err)
	assert.True(t, next > 0)

	assert.Equal(t, now+1, next)
}
