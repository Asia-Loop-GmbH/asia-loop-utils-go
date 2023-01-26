package normalizer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v7/pkg/normalizer"
)

func TestEmail(t *testing.T) {
	assert.Equal(t, "lenamtruong@gmail.com", normalizer.Email(context.Background(), "LeNamtrUong@gmail.com"))
}

func TestPhoneNumber(t *testing.T) {
	assert.Equal(t, "+49 170 1234567", normalizer.PhoneNumber(context.Background(), "1701234567"))
}

func TestName(t *testing.T) {
	assert.Equal(t, "Le Nam-Truong Nhung", normalizer.Name(context.Background(), "  le     nam-truong     nhung  "))
}
