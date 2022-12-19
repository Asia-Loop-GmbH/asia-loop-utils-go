package api

import (
	"github.com/asia-loop-gmbh/asia-loop-utils-go/pkg/db"
)

type GroupDetails struct {
	Group  db.Group       `json:"group"`
	Orders []OrderDetails `json:"orders"`
}

type CreateGroupRequest struct {
	Store string `json:"store"`
}