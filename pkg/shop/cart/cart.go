package cart

import (
	"context"
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

func toOrder(ctx context.Context, shoppingCart *db.Cart, coupon *db.Coupon, products []db.Product, taxes []db.Tax) (*db.Order, error) {
	log := logger.FromContext(ctx)
	log.Infof("Calculate cart")

	sTotal := decimal.Zero
	sTax := decimal.Zero
	sNet := decimal.Zero
	sSaving := decimal.Zero

	items := lo.Map(shoppingCart.Items, func(item db.CartItem, index int) db.OrderItem {
		product, ok := lo.Find(products, func(p db.Product) bool {
			return p.ID.Hex() == item.ProductID
		})
		if !ok {
			return db.OrderItem{} // TODO: we don't expect that this happens, so just ignore this case for now
		}
		tax, ok := lo.Find(taxes, func(tax db.Tax) bool {
			return tax.Name == db.TaxClassTakeaway
		})
		if !ok {
			return db.OrderItem{}
		}

		itemPrice := decimal.RequireFromString(product.GetPrice(shoppingCart.StoreKey, item.Options))
		saving := decimal.Zero
		if shoppingCart.IsPickup && !product.IsGiftCard {
			originalPrice := itemPrice.Add(decimal.Zero)
			itemPrice = itemPrice.Mul(decimal.NewFromFloat(0.8)).Round(2)
			saving = originalPrice.Sub(itemPrice).Mul(decimal.NewFromInt(int64(item.Amount)))
			sSaving = sSaving.Add(saving)
		}
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

		return db.OrderItem{
			CartItem:   item,
			SKU:        product.SKU,
			Name:       product.Name,
			Categories: product.Categories,
			UnitPrice:  itemPrice.StringFixed(2), // TODO: we must improve this in the future, when we support store specific prices
			Total:      totalPrice.StringFixed(2),
			Tax:        totalTax.StringFixed(2),
			Net:        totalNet.StringFixed(2),
			Saving:     saving.StringFixed(2),
			TaxClass:   db.TaxClassTakeaway,
			IsGiftCard: product.IsGiftCard,
		}
	})
	if coupon != nil {
		couponAmount := decimal.RequireFromString(coupon.Available())
		couponAmount = decimal.Min(couponAmount, sTotal)
		couponTax := couponAmount.Div(decimal.NewFromFloat(1.07)).Mul(decimal.NewFromFloat(0.07)).Round(2)
		couponNet := couponAmount.Sub(couponTax)

		sTotal = sTotal.Sub(couponAmount)
		sNet = sNet.Sub(couponNet)
		sTax = sTax.Sub(couponTax)

		items = append(items, db.OrderItem{
			CartItem: db.CartItem{
				Amount: 1,
			},
			SKU:          db.CouponSKU,
			Name:         "Gutschein",
			Categories:   nil,
			UnitPrice:    couponAmount.Neg().StringFixed(2),
			Total:        couponAmount.Neg().StringFixed(2),
			Tax:          couponTax.Neg().StringFixed(2),
			Net:          couponNet.Neg().StringFixed(2),
			Saving:       "0.00",
			TaxClass:     db.TaxClassTakeaway,
			IsGiftCard:   false,
			GiftCardCode: nil,
		})
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

	summary := db.OrderSummary{
		Total: db.TotalSummary{
			Value:  sTotal.StringFixed(2),
			Values: map[string]string{db.TaxClassTakeaway: sTotal.StringFixed(2)},
		},
		Tax: db.TotalSummary{
			Value:  sTax.StringFixed(2),
			Values: map[string]string{db.TaxClassTakeaway: sTax.StringFixed(2)},
		},
		Net: db.TotalSummary{
			Value:  sNet.StringFixed(2),
			Values: map[string]string{db.TaxClassTakeaway: sNet.StringFixed(2)},
		},
		Saving: sSaving.StringFixed(2),
	}
	if tip != nil {
		summary.Total.Value = sTotal.Add(*tip).StringFixed(2)
		summary.Total.Values[db.TaxClassZero] = tip.StringFixed(2)

		summary.Net.Values[db.TaxClassZero] = tip.StringFixed(2)
		summary.Net.Value = sNet.Add(*tip).StringFixed(2)
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
