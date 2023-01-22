package servicegooglemaps_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v6/pkg/servicegooglemaps"
	commoncontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
)

func TestResolveAddress(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.TODO(), commoncontext.FieldStage, "dev")
	result, err := servicegooglemaps.ResolveAddress(ctx, "rudolf-schwarz platz 1 frankfurt")
	expected := &servicegooglemaps.ResolveAddressResult{
		StreetNumber:     "1",
		Street:           "Rudolf-Schwarz-Platz",
		City:             "Frankfurt am Main",
		Postcode:         "60438",
		State:            "Hessen",
		FormattedAddress: "Rudolf-Schwarz-Platz 1, 60438 Frankfurt am Main, Deutschland",
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
