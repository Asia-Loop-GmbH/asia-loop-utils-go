package api

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetRevenueRequest struct {
	From    time.Time          `json:"from"`
	To      time.Time          `json:"to"`
	StoreID primitive.ObjectID `json:"storeId"`
}

type Revenue struct {
	Request GetRevenueRequest `json:"request"`
	Items   []RevenueItem     `json:"items"`
}

type RevenueItem struct {
	Text  string               `json:"type"`
	Total primitive.Decimal128 `json:"total"`
}
