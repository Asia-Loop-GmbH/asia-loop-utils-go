package cart

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v8/pkg/shop/db"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

var (
	deliveryDiscountExcludedProductCategories = []string{"drink"}
)

func ToOrder(ctx context.Context, shoppingCart *db.Cart) (*db.Order, error) {
	log := logger.FromContext(ctx)
	log.Infof("Calculate public cart")

	colProducts, err := db.CollectionProducts(ctx)
	if err != nil {
		log.Errorf("Failed to init db collection: %s", err)
		return nil, errors.Wrap(err, "failed to init db collection")
	}
	findProducts, err := colProducts.Find(ctx, bson.M{}) // TODO: search only related products
	if err != nil {
		log.Errorf("Failed to find products: %s", err)
		return nil, errors.Wrap(err, "failed to find products")
	}
	products := make([]db.Product, 0)
	err = findProducts.All(ctx, &products)
	if err != nil {
		log.Errorf("Failed to decode products: %s", err)
		return nil, errors.Wrap(err, "failed to decode products")
	}

	colTaxes, err := db.CollectionTaxes(ctx)
	if err != nil {
		log.Errorf("Failed to init db collection: %s", err)
		return nil, errors.Wrap(err, "failed to init db collection")
	}
	findTaxes, err := colTaxes.Find(ctx, bson.M{})
	if err != nil {
		log.Errorf("Failed to find taxes: %s", err)
		return nil, errors.Wrap(err, "failed to find taxes")
	}
	taxes := make([]db.Tax, 0)
	err = findTaxes.All(ctx, &taxes)
	if err != nil {
		log.Errorf("Failed to decode taxes: %s", err)
		return nil, errors.Wrap(err, "failed to decode taxes")
	}

	coupon := new(db.Coupon)
	if shoppingCart.CouponCode != nil {
		code := strings.ToUpper(*shoppingCart.CouponCode)

		colCoupons, err := db.CollectionCoupons(ctx)
		if err != nil {
			log.Errorf("Failed to init db collection: %s", err)
			return nil, errors.Wrap(err, "failed to init db collection")
		}
		findCoupon := colCoupons.FindOne(ctx, bson.M{"code": code, "disabled": false})
		err = findCoupon.Decode(coupon)
		if err == mongo.ErrNoDocuments {
			log.Errorf("Coupon not found: %s", code)
			coupon = nil
		} else if err != nil {
			log.Errorf("Failed to decode coupon: %s", err)
			return nil, errors.Wrap(err, "failed to decode coupon")
		}
	} else {
		coupon = nil
	}

	return toOrder(ctx, shoppingCart, coupon, products, taxes)
}

type summaryTotal struct {
	Total  decimal.Decimal
	Tax    decimal.Decimal
	Net    decimal.Decimal
	Saving decimal.Decimal
}

func toOrder(ctx context.Context, shoppingCart *db.Cart, coupon *db.Coupon, products []db.Product, taxes []db.Tax) (*db.Order, error) {
	log := logger.FromContext(ctx)
	log.Infof("Calculate cart")

	takeawaySummary := summaryTotal{
		Total:  decimal.Zero,
		Tax:    decimal.Zero,
		Net:    decimal.Zero,
		Saving: decimal.Zero,
	}
	standardSummary := summaryTotal{
		Total:  decimal.Zero,
		Tax:    decimal.Zero,
		Net:    decimal.Zero,
		Saving: decimal.Zero,
	}
	zeroSummary := summaryTotal{
		Total:  decimal.Zero,
		Tax:    decimal.Zero,
		Net:    decimal.Zero,
		Saving: decimal.Zero,
	}
	summaryTotals := map[string]*summaryTotal{
		db.TaxClassTakeaway: &takeawaySummary,
		db.TaxClassStandard: &standardSummary,
		db.TaxClassZero:     &zeroSummary,
	}

	items := lo.Map(shoppingCart.Items, func(item db.CartItem, index int) db.OrderItem {
		product, ok := lo.Find(products, func(p db.Product) bool {
			return p.ID.Hex() == item.ProductID
		})
		if !ok {
			return db.OrderItem{} // TODO: we don't expect that this happens, so just ignore this case for now
		}
		tax, ok := lo.Find(taxes, func(tax db.Tax) bool {
			return tax.Name == product.TaxClassTakeAway // we don't sell in store, so we can only use this tax
		})
		if !ok {
			return db.OrderItem{}
		}
		summary, ok := summaryTotals[tax.Name]
		if !ok {
			return db.OrderItem{}
		}

		itemPrice := decimal.RequireFromString(product.GetPrice(shoppingCart.StoreKey, item.Options))
		saving := decimal.Zero
		if shoppingCart.IsPickup && !product.IsGiftCard && len(lo.Intersect(deliveryDiscountExcludedProductCategories, product.Categories)) == 0 {
			originalPrice := itemPrice.Add(decimal.Zero)
			itemPrice = itemPrice.Mul(decimal.NewFromFloat(0.8)).Round(2)
			saving = originalPrice.Sub(itemPrice).Mul(decimal.NewFromInt(int64(item.Amount)))
			summary.Saving = summary.Saving.Add(saving)
		}
		taxRate := decimal.RequireFromString(tax.Rate)
		taxPrice := itemPrice.Div(decimal.NewFromInt(1).Add(taxRate)).Mul(taxRate).Round(2)
		netPrice := itemPrice.Sub(taxPrice)

		amount := decimal.NewFromInt(int64(item.Amount))
		totalPrice := itemPrice.Mul(amount)
		totalNet := netPrice.Mul(amount)
		totalTax := taxPrice.Mul(amount)

		summary.Total = summary.Total.Add(totalPrice)
		summary.Net = summary.Net.Add(totalNet)
		summary.Tax = summary.Tax.Add(totalTax)

		return db.OrderItem{
			CartItem:   item,
			SKU:        product.SKU,
			Name:       product.Name,
			Categories: product.Categories,
			UnitPrice:  itemPrice.StringFixed(2),
			Total:      totalPrice.StringFixed(2),
			Tax:        totalTax.StringFixed(2),
			Net:        totalNet.StringFixed(2),
			Saving:     saving.StringFixed(2),
			TaxClass:   tax.Name,
			IsGiftCard: product.IsGiftCard,
		}
	})

	appliedCouponItems := generateCouponItems(ctx, coupon, &takeawaySummary, &standardSummary)

	for _, item := range appliedCouponItems {
		switch item.TaxClass {
		case db.TaxClassTakeaway:
			takeawaySummary.Total = takeawaySummary.Total.Add(decimal.RequireFromString(item.Total))
			takeawaySummary.Tax = takeawaySummary.Tax.Add(decimal.RequireFromString(item.Tax))
			takeawaySummary.Net = takeawaySummary.Net.Add(decimal.RequireFromString(item.Net))
		case db.TaxClassStandard:
			standardSummary.Total = standardSummary.Total.Add(decimal.RequireFromString(item.Total))
			standardSummary.Tax = standardSummary.Tax.Add(decimal.RequireFromString(item.Tax))
			standardSummary.Net = standardSummary.Net.Add(decimal.RequireFromString(item.Net))
		default:
			log.Errorf("Unsupported coupon tax [%s]", item.TaxClass)
			return nil, fmt.Errorf("unsupported coupon tax [%s]", item.TaxClass)
		}
		items = append(items, item)
	}

	var tip *decimal.Decimal
	if shoppingCart.Tip != nil {
		t, err := decimal.NewFromString(*shoppingCart.Tip)
		if err != nil {
			log.Errorf("Invalid trip: %s", *shoppingCart.Tip)
		} else {
			tip = &t

			items = append(items, db.OrderItem{
				CartItem: db.CartItem{
					Amount: 1,
				},
				SKU:          db.TipSKU,
				Name:         "Trinkgeld",
				Categories:   nil,
				UnitPrice:    tip.StringFixed(2),
				Total:        tip.StringFixed(2),
				Tax:          "0.00",
				Net:          tip.StringFixed(2),
				Saving:       "0.00",
				TaxClass:     db.TaxClassZero,
				IsGiftCard:   false,
				GiftCardCode: nil,
			})
		}
	}

	sTotal := takeawaySummary.Total.Add(standardSummary.Total).Add(zeroSummary.Total)
	sTax := takeawaySummary.Tax.Add(standardSummary.Tax).Add(zeroSummary.Tax)
	sNet := takeawaySummary.Net.Add(standardSummary.Net).Add(zeroSummary.Net)
	sSaving := takeawaySummary.Saving.Add(standardSummary.Saving).Add(zeroSummary.Saving)

	totalValues := map[string]string{}
	totalValues[db.TaxClassTakeaway] = takeawaySummary.Total.StringFixed(2)
	if !standardSummary.Total.IsZero() {
		totalValues[db.TaxClassStandard] = standardSummary.Total.StringFixed(2)
	}
	if !zeroSummary.Total.IsZero() {
		totalValues[db.TaxClassZero] = zeroSummary.Total.StringFixed(2)
	}

	taxValues := map[string]string{}
	taxValues[db.TaxClassTakeaway] = takeawaySummary.Tax.StringFixed(2)
	if !standardSummary.Tax.IsZero() {
		taxValues[db.TaxClassStandard] = standardSummary.Tax.StringFixed(2)
	}
	if !zeroSummary.Tax.IsZero() {
		taxValues[db.TaxClassZero] = zeroSummary.Tax.StringFixed(2)
	}

	netValues := map[string]string{}
	netValues[db.TaxClassTakeaway] = takeawaySummary.Net.StringFixed(2)
	if !standardSummary.Net.IsZero() {
		netValues[db.TaxClassStandard] = standardSummary.Net.StringFixed(2)
	}
	if !zeroSummary.Net.IsZero() {
		netValues[db.TaxClassZero] = zeroSummary.Net.StringFixed(2)
	}

	summary := db.OrderSummary{
		Total: db.TotalSummary{
			Value:  sTotal.StringFixed(2),
			Values: totalValues,
		},
		Tax: db.TotalSummary{
			Value:  sTax.StringFixed(2),
			Values: taxValues,
		},
		Net: db.TotalSummary{
			Value:  sNet.StringFixed(2),
			Values: netValues,
		},
		Saving: sSaving.StringFixed(2),
	}
	if tip != nil {
		summary.Total.Value = sTotal.Add(*tip).StringFixed(2)
		if currentTotal, found := summary.Total.Values[db.TaxClassZero]; found {
			summary.Total.Values[db.TaxClassZero] = tip.Add(decimal.RequireFromString(currentTotal)).StringFixed(2)
		} else {
			summary.Total.Values[db.TaxClassZero] = tip.StringFixed(2)
		}

		summary.Net.Value = sNet.Add(*tip).StringFixed(2)
		if currentNet, found := summary.Net.Values[db.TaxClassZero]; found {
			summary.Net.Values[db.TaxClassZero] = tip.Add(decimal.RequireFromString(currentNet)).StringFixed(2)
		} else {
			summary.Net.Values[db.TaxClassZero] = tip.StringFixed(2)
		}
	}

	var payment *db.Payment
	last, err := lo.Last(shoppingCart.Payments)
	if err == nil {
		finalTotal := decimal.RequireFromString(summary.Total.Value)
		sameAmount := last.Session.Amount.Value == finalTotal.Mul(decimal.NewFromInt(100)).IntPart()
		notExpired := last.Session.ExpiresAt.After(time.Now())
		if sameAmount && notExpired {
			payment = &last
		}
	}

	return &db.Order{
		ID:         shoppingCart.ID,
		User:       shoppingCart.User,
		CouponCode: shoppingCart.CouponCode,
		Tip:        shoppingCart.Tip,
		StoreKey:   shoppingCart.StoreKey,
		Checkout:   shoppingCart.Checkout,
		IsPickup:   shoppingCart.IsPickup,
		Paid:       shoppingCart.Paid,
		Items:      items,
		Summary:    summary,
		Payment:    payment,
		Secret:     shoppingCart.Secret,
		CreatedAt:  shoppingCart.CreatedAt,
		UpdatedAt:  shoppingCart.UpdatedAt,
	}, nil
}
