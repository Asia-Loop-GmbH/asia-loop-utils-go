package globalconfig_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v5/pkg/globalconfig"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
)

func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.Background()
	log := logger.FromContext(ctx)

	cfg, err := globalconfig.Get(ctx)
	assert.NoError(t, err)
	log.Infof("%v", cfg)
}
