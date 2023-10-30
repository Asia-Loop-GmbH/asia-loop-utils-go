package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mongodb"
)

const colOrders = "orders"

func CollectionOrders(ctx context.Context) (*mongo.Collection, error) {
	database, err := secretsmanager.GetParameter(ctx, "/mongo/database")
	if err != nil {
		return nil, err
	}
	return mongodb.Collection(ctx, database, colOrders)
}

type OrderShippingMethod string

const (
	OrderShippingMethodFreeShipping OrderShippingMethod = "free_shipping"
	OrderShippingMethodLocalPickup  OrderShippingMethod = "local_pickup"
	OrderShippingMethodNoShipping   OrderShippingMethod = "no_shipping"
)

type OrderStatus string

const (
	OrderStatusPending        OrderStatus = "PENDING"
	OrderStatusNew            OrderStatus = "NEW"
	OrderStatusNotPossible    OrderStatus = "NOT_POSSIBLE"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
	OrderStatusAddedToGroup   OrderStatus = "ADDED_TO_GROUP"
	OrderStatusGroupFinalized OrderStatus = "GROUP_FINALIZED"
	OrderStatusDelivered      OrderStatus = "DELIVERED"
)

type OrderItemExtra struct {
	Key   string `bson:"key" json:"key"`
	Value string `bson:"value" json:"value"`
}

type OrderItem struct {
	Name       string           `bson:"name" json:"name"`
	Quantity   int              `bson:"quantity" json:"quantity"`
	SKU        string           `bson:"sku" json:"sku"`
	Categories []string         `bson:"categories" json:"categories"`
	Extra      []OrderItemExtra `bson:"extra" json:"extra"`
	Net        string           `bson:"net" json:"net"`
	Tax        string           `bson:"tax" json:"tax"`
	TaxClass   TaxClass         `bson:"taxClass" json:"taxClass"`
	TaxRate    string           `bson:"taxRate" json:"taxRate"`
	Total      string           `bson:"total" json:"total"`
}

var OrderObjectIDs = []string{
	"customer", "store",
}

type Order struct {
	ID                   primitive.ObjectID  `bson:"_id" json:"id"`
	OrderID              string              `bson:"orderId" json:"orderId"`
	OrderNumber          string              `bson:"orderNumber" json:"orderNumber"`
	InvoiceNumber        string              `bson:"invoiceNumber" json:"invoiceNumber"`
	ValidAddress         bool                `bson:"validAddress" json:"validAddress"`
	Address              string              `bson:"address" json:"address"`
	AddressLine2         string              `bson:"addressLine2" json:"addressLine2"`
	Email                string              `bson:"email" json:"email"`
	Telephone            string              `bson:"telephone" json:"telephone"`
	CustomerNote         string              `bson:"customerNote" json:"customerNote"`
	DeliveryDate         string              `bson:"deliveryDate" json:"deliveryDate"`
	DeliveryTime         string              `bson:"deliveryTime" json:"deliveryTime"`
	ShippingMethod       OrderShippingMethod `bson:"shippingMethod" json:"shippingMethod"`
	Customer             primitive.ObjectID  `bson:"customer" json:"customer"`
	Status               OrderStatus         `bson:"status" json:"status"`
	Secret               string              `bson:"secret" json:"secret"`
	Items                []OrderItem         `bson:"items" json:"items"`
	Net                  string              `bson:"net" json:"net"`
	Tax                  string              `bson:"tax" json:"tax"`
	Total                string              `bson:"total" json:"total"`
	PaidTotal            string              `bson:"paidTotal" json:"paidTotal"`
	CouponCode           string              `bson:"couponCode" json:"couponCode"`
	AppliedCouponNet     string              `bson:"appliedCouponNet" json:"appliedCouponNet"`
	AppliedCouponTax     string              `bson:"appliedCouponTax" json:"appliedCouponTax"`
	AppliedCouponTotal   string              `bson:"appliedCouponTotal" json:"appliedCouponTotal"`
	AppliedCoupon7Net    string              `bson:"appliedCoupon7Net" json:"appliedCoupon7Net"`
	AppliedCoupon7Tax    string              `bson:"appliedCoupon7Tax" json:"appliedCoupon7Tax"`
	AppliedCoupon7Total  string              `bson:"appliedCoupon7Total" json:"appliedCoupon7Total"`
	AppliedCoupon19Net   string              `bson:"appliedCoupon19Net" json:"appliedCoupon19Net"`
	AppliedCoupon19Tax   string              `bson:"appliedCoupon19Tax" json:"appliedCoupon19Tax"`
	AppliedCoupon19Total string              `bson:"appliedCoupon19Total" json:"appliedCoupon19Total"`
	Tip                  string              `bson:"tip" json:"tip"`
	Printed              bool                `bson:"printed" json:"printed"`
	LocalPickupNotified  bool                `bson:"localPickupNotified" json:"localPickupNotified"`
	ReadyIn              string              `bson:"readyIn" json:"readyIn"`
	Store                primitive.ObjectID  `bson:"store" json:"store"`
	CompanyKey           string              `bson:"companyKey" json:"companyKey"`
	PaymentEvents        []interface{}       `bson:"paymentEvents" json:"paymentEvents"`
	LastSessionId        string              `bson:"lastSessionId" json:"lastSessionId"` // for corporate page
	NewBoxes             int                 `bson:"newBoxes" json:"newBoxes"`
	ReturnBoxes          int                 `bson:"returnBoxes" json:"returnBoxes"`
	CreatedAt            time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt            time.Time           `bson:"updatedAt" json:"updatedAt"`
}
