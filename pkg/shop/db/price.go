package db

import (
	"github.com/shopspring/decimal"
)

type CustomizablePrice struct {
	Value        string            `bson:"value" json:"value"`
	CustomValues map[string]string `bson:"customValues" json:"customValues"`
}

func (p *CustomizablePrice) GetPrice(key string) string {
	custom, ok := p.CustomValues[key]
	if ok && decimal.RequireFromString(custom).IsPositive() {
		return custom
	}
	return p.Value
}
