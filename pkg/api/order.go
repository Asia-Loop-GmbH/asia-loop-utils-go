package api

import (
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v9/pkg/db"
)

type OrderDetails struct {
	Order    db.Order    `json:"order"`
	Customer db.Customer `json:"customer"`
}

type SearchOrderRequest struct {
	Text  *string `json:"text"`
	Limit *int64  `json:"limit"`
}
