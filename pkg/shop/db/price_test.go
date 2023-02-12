package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomizablePrice_GetPrice(t *testing.T) {
	assert.Equal(t, "1.23", (&CustomizablePrice{
		Value:        "1.23",
		CustomValues: nil,
	}).GetPrice("FOO"))

	assert.Equal(t, "4.56", (&CustomizablePrice{
		Value: "1.23",
		CustomValues: map[string]string{
			"FOO": "4.56",
		},
	}).GetPrice("FOO"))

	assert.Equal(t, "1.23", (&CustomizablePrice{
		Value: "1.23",
		CustomValues: map[string]string{
			"FOO": "0.00",
		},
	}).GetPrice("FOO"))

	assert.Equal(t, "1.23", (&CustomizablePrice{
		Value: "1.23",
		CustomValues: map[string]string{
			"FOO": "0",
		},
	}).GetPrice("FOO"))
}
