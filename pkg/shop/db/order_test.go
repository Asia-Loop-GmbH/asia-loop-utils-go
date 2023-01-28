package db_test

import (
	"testing"

	"github.com/adyen/adyen-go-api-library/v6/src/notification"
	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v7/pkg/shop/db"
)

func TestOrder_GetPaidEvent(t *testing.T) {
	order := db.Order{
		Payment: &db.Payment{
			Events: []notification.NotificationRequestItem{
				{
					EventCode: "foo",
					Success:   "bar",
				},
				{
					EventCode: "AUTHORISATION",
					Success:   "bar",
				},
				{
					EventCode: "AUTHORISATION",
					Success:   "true",
				},
				{
					EventCode: "foo",
					Success:   "true",
				},
			},
		},
	}

	assert.Equal(t, notification.NotificationRequestItem{
		EventCode: "AUTHORISATION", Success: "true",
	}, order.GetPaidEvent())
}
