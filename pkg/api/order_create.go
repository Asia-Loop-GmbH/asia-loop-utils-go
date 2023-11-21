package api

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v9/pkg/db"
)

type CreateOrderRequest struct {
	FirstName            string                 `json:"firstName"`
	LastName             string                 `json:"lastName"`
	Telephone            string                 `json:"telephone"`
	Email                string                 `json:"email"`
	AddressLine1         string                 `json:"addressLine1"`
	AddressLine2         string                 `json:"addressLine2"`
	Postcode             string                 `json:"postcode"`
	City                 string                 `json:"city"`
	ShippingMethod       db.OrderShippingMethod `json:"shippingMethod"`
	Items                []db.OrderItem         `json:"items"`
	Net                  string                 `json:"net"`
	Tax                  string                 `json:"tax"`
	Total                string                 `json:"total"`
	CouponCode           string                 `json:"couponCode"`
	AppliedCouponNet     string                 `json:"appliedCouponNet"`
	AppliedCouponTax     string                 `json:"appliedCouponTax"`
	AppliedCouponTotal   string                 `json:"appliedCouponTotal"`
	AppliedCoupon7Net    string                 `json:"appliedCoupon7Net"`
	AppliedCoupon7Tax    string                 `json:"appliedCoupon7Tax"`
	AppliedCoupon7Total  string                 `json:"appliedCoupon7Total"`
	AppliedCoupon19Net   string                 `json:"appliedCoupon19Net"`
	AppliedCoupon19Tax   string                 `json:"appliedCoupon19Tax"`
	AppliedCoupon19Total string                 `json:"appliedCoupon19Total"`
	Store                primitive.ObjectID     `json:"store"`
	CompanyKey           string                 `json:"companyKey"`
}

type CreateOrderAddressOptions struct {
	FirstName    string
	LastName     string
	Telephone    string
	Email        string
	AddressLine1 string
	AddressLine2 string
	Postcode     string
	City         string
}

type CreateOrderOrderOptions struct {
	OrderID              string
	OrderNumber          string
	InvoiceNumber        string
	Items                []db.OrderItem
	DeliveryDate         string
	DeliveryTime         string
	Net                  string
	Total                string
	Tax                  string
	Tip                  string
	CouponCode           string
	AppliedCouponNet     string
	AppliedCouponTax     string
	AppliedCouponTotal   string
	AppliedCoupon7Net    string
	AppliedCoupon7Tax    string
	AppliedCoupon7Total  string
	AppliedCoupon19Net   string
	AppliedCoupon19Tax   string
	AppliedCoupon19Total string
	ShippingMethod       db.OrderShippingMethod
	CustomerNote         string
	Status               db.OrderStatus
	Store                primitive.ObjectID
	CompanyKey           string
}
