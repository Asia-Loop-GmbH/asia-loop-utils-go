package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v7/pkg/shop/db"
)

func TestCoupon_Available_Disabled(t *testing.T) {
	c := &db.Coupon{
		Total:    "10.00",
		Disabled: true,
	}
	assert.Equal(t, "0.00", c.Available())
}

func TestCoupon_Available(t *testing.T) {
	c := &db.Coupon{
		Total:    "10.00",
		Disabled: false,
		Usage: []db.CouponUsage{
			{
				Total: "2.00",
			},
			{
				Total: "3.00",
			},
		},
	}
	assert.Equal(t, "5.00", c.Available())
}
