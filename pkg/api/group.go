package api

import (
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v3/pkg/db"
)

type GroupDetails struct {
	Group  db.Group       `json:"group"`
	Orders []OrderDetails `json:"orders"`
}

type CreateGroupRequest struct {
	Store string `json:"store"`
}
