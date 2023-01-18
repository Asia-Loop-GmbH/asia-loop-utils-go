package cart

import (
	"context"
	"time"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v2/pkg/shop/db"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
)

type PublicCart struct {
	ID        primitive.ObjectID `json:"id"`
	IsPickup  bool               `json:"isPickup"`
	Items     []PublicCartItem   `json:"items"`
	Summary   PublicCartSummary  `json:"summary"`
	Secret    string             `json:"secret"`
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
	SKU        string   `json:"sku"`
	Categories []string `json:"categories"`
	UnitPrice  string   `json:"unitPrice"`
	Total      string   `json:"total"`
	Tax        string   `json:"tax"`
	Net        string   `json:"net"`
	TaxClass   string   `json:"taxClass"`
}

func Calculate(ctx context.Context, c *db.Cart, products []db.Product, taxes []db.Tax) (*PublicCart, error) {
	log := logger.FromContext(ctx)
	log.Infof("Calculate cart")

	sTotal := decimal.Zero
	sTax := decimal.Zero
	sNet := decimal.Zero
	sTaxes := map[string]decimal.Decimal{
		"takeaway": decimal.Zero,
	}

	items := lo.Map(c.Items, func(item db.CartItem, index int) PublicCartItem {
		product, ok := lo.Find(products, func(p db.Product) bool {
			return p.ID.String() == item.ProductID
		})
		if !ok {
			return PublicCartItem{} // TODO: we don't expect that this happens, so just ignore this case for now
		}
		tax, ok := lo.Find(taxes, func(tax db.Tax) bool {
			return tax.Name == "takeaway"
		})
		if !ok {
			return PublicCartItem{}
		}

		itemPrice := decimal.RequireFromString(product.Price.Value)
		taxRate := decimal.RequireFromString(tax.Rate)
		taxPrice := itemPrice.Div(decimal.NewFromInt(1).Add(taxRate)).Mul(taxRate).Round(2)
		netPrice := itemPrice.Sub(taxPrice)

		amount := decimal.NewFromInt(int64(item.Amount))
		totalPrice := itemPrice.Mul(amount)
		totalNet := netPrice.Mul(amount)
		totalTax := taxPrice.Mul(amount)

		sTotal = sTotal.Add(totalPrice)
		sNet = sNet.Add(totalNet)
		sTax = sTax.Add(totalTax)
		sTaxes["takeaway"] = sTaxes["takeaway"].Add(sTax) // TODO: this must be improved in the future, when we support multiple taxes

		return PublicCartItem{
			CartItem:   item,
			SKU:        product.SKU,
			Categories: product.Categories,
			UnitPrice:  product.Price.Value, // TODO: we must improve this in the future, when we support store specific prices
			Total:      totalPrice.StringFixed(2),
			Tax:        totalTax.StringFixed(2),
			Net:        totalNet.StringFixed(2),
			TaxClass:   "takeaway",
		}
	})

	return &PublicCart{
		ID:       c.ID,
		IsPickup: c.IsPickup,
		Items:    items,
		Summary: PublicCartSummary{
			Total:    sTotal.StringFixed(2),
			TotalTax: sTax.StringFixed(2),
			TotalNet: sNet.StringFixed(2),
			Taxes: map[string]string{
				"takeaway": sTaxes["takeaway"].StringFixed(2),
			},
		},
		Secret:    c.Secret,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}, nil
}
