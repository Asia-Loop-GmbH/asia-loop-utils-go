package cart

import (
	"context"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v2/pkg/shop/db"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
)

type PublicCart struct {
	ID        primitive.ObjectID `json:"id"`
	Items     []PublicCartItem   `json:"items"`
	Summary   PublicCartSummary  `json:"summary"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type PublicCartSummary struct {
	Total    string            `json:"total"`
	TotalTax string            `json:"totalTax"`
	TotalNet string            `json:"net"`
	Taxes    map[string]string `json:"taxes"`
}

type PublicCartItem struct {
	db.CartItem
	UnitPrice string `json:"unitPrice"`
	Total     string `json:"total"`
	Tax       string `json:"tax"`
	Net       string `json:"net"`
	TaxClass  string `json:"taxClass"`
}

func Calculate(ctx context.Context, c *db.Cart) (*PublicCart, error) {
	log := logger.FromContext(ctx)
	log.Infof("Calculate cart")

	return &PublicCart{
		ID: c.ID,
		Items: lo.Map(c.Items, func(item db.CartItem, index int) PublicCartItem {
			return PublicCartItem{
				CartItem:  item,
				UnitPrice: "0.00",
				Total:     "0.00",
				Tax:       "0.00",
				Net:       "0.00",
				TaxClass:  "takeaway",
			}
		}),
		Summary: PublicCartSummary{
			Total:    "0.00",
			TotalTax: "0.00",
			TotalNet: "0.00",
			Taxes: map[string]string{
				"MwSt. 7%": "0.00",
			},
		},
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}, nil
}
