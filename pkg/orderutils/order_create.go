package orderutils

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v8/pkg/api"
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v8/pkg/db"
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v8/pkg/normalizer"
	"github.com/asia-loop-gmbh/asia-loop-utils-go/v8/pkg/servicegooglemaps"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/random"
)

func CreateOrder(ctx context.Context, orderOptions *api.CreateOrderOrderOptions,
	addressOption *api.CreateOrderAddressOptions) (*db.Order, error) {
	log := logger.FromContext(ctx)

	log.Infof("create order: %s", orderOptions.OrderID)

	firstName := normalizer.Name(ctx, addressOption.FirstName)
	lastName := normalizer.Name(ctx, addressOption.LastName)
	telephone := normalizer.PhoneNumber(ctx, addressOption.Telephone)
	email := normalizer.Email(ctx, addressOption.Email)

	var addressLine1 string
	addressLine2 := addressOption.AddressLine2
	var postcode string
	var city string
	var formattedAddress string
	var validAddress bool

	inputAddress := fmt.Sprintf(
		"%s, %s %s",
		addressOption.AddressLine1,
		addressOption.Postcode,
		addressOption.City,
	)

	resolveAddressResult, err := servicegooglemaps.ResolveAddress(ctx, inputAddress)
	if err != nil {
		log.Warnf("failed to resolve address: %s", inputAddress)
		addressLine1 = addressOption.AddressLine1
		postcode = addressOption.Postcode
		city = addressOption.City
		formattedAddress = inputAddress
		validAddress = false
	} else {
		addressLine1 = fmt.Sprintf("%s %s", resolveAddressResult.Street, resolveAddressResult.StreetNumber)
		postcode = resolveAddressResult.Postcode
		city = resolveAddressResult.City
		formattedAddress = resolveAddressResult.FormattedAddress
		validAddress = true
	}

	collectionCustomer, err := db.CollectionCustomers(ctx)
	if err != nil {
		return nil, err
	}
	findCustomer := collectionCustomer.FindOne(ctx, bson.M{
		"addressLine1": addressLine1,
		"addressLine2": addressLine2,
		"postcode":     postcode,
		"city":         city,
		"firstName":    firstName,
		"lastName":     lastName,
	})
	customer := new(db.Customer)
	err = findCustomer.Decode(customer)
	if err != nil {
		log.Infof("failed to find customer: %s %s (%s)", firstName, lastName, formattedAddress)
		customerRef := fmt.Sprintf(
			"%s%s",
			random.String(2, lo.UpperCaseLettersCharset),
			random.String(6, lo.NumbersCharset),
		)
		customerCreatedAt := time.Now()
		newCustomer := db.Customer{
			ID:           primitive.NewObjectID(),
			FirstName:    firstName,
			LastName:     lastName,
			CustomerRef:  customerRef,
			Boxes:        []int{},
			AddressLine1: addressLine1,
			AddressLine2: addressLine2,
			Postcode:     postcode,
			City:         city,
			Telephone:    telephone,
			Email:        email,
			CreatedAt:    customerCreatedAt,
			UpdatedAt:    customerCreatedAt,
		}
		if _, err := collectionCustomer.InsertOne(ctx, newCustomer); err != nil {
			return nil, err
		}
		log.Infof("new customer created: %s %s (%s)", firstName, lastName, formattedAddress)
		customer = &newCustomer
	}

	collectionOrder, err := db.CollectionOrders(ctx)
	if err != nil {
		return nil, err
	}
	secret := random.String(32, lo.AlphanumericCharset)
	orderCreatedAt := time.Now()
	newOrder := db.Order{
		ID:                  primitive.NewObjectID(),
		Status:              orderOptions.Status,
		OrderID:             orderOptions.OrderID,
		OrderNumber:         orderOptions.OrderNumber,
		InvoiceNumber:       orderOptions.InvoiceNumber,
		ValidAddress:        validAddress,
		Address:             formattedAddress,
		AddressLine2:        addressLine2,
		Email:               email,
		Telephone:           telephone,
		CustomerNote:        orderOptions.CustomerNote,
		DeliveryDate:        orderOptions.DeliveryDate,
		DeliveryTime:        orderOptions.DeliveryTime,
		ShippingMethod:      orderOptions.ShippingMethod,
		Customer:            customer.ID,
		Items:               orderOptions.Items,
		Tip:                 orderOptions.Tip,
		Total:               orderOptions.Total,
		Tax:                 orderOptions.Tax,
		Net:                 orderOptions.Net,
		CouponCode:          orderOptions.CouponCode,
		AppliedCouponTax:    orderOptions.AppliedCouponTax,
		AppliedCouponTotal:  orderOptions.AppliedCouponTotal,
		AppliedCouponNet:    orderOptions.AppliedCouponNet,
		Store:               orderOptions.Store,
		CompanyKey:          orderOptions.CompanyKey,
		Secret:              secret,
		Printed:             false,
		LocalPickupNotified: false,
		PaymentEvents:       make([]interface{}, 0),
		CreatedAt:           orderCreatedAt,
		UpdatedAt:           orderCreatedAt,
	}

	if _, err := collectionOrder.InsertOne(ctx, newOrder); err != nil {
		return nil, err
	}

	log.Infof("order created: %s", orderOptions.OrderID)
	customerUpdatedAt := time.Now()
	if customer.Telephone == "" {
		_, err := collectionCustomer.UpdateByID(ctx, customer.ID, bson.M{
			"$set": bson.M{
				"telephone": telephone,
				"updatedAt": customerUpdatedAt,
			},
		})
		if err != nil {
			log.Errorf("failed to update customer telphone to %s: %s", telephone, err)
		} else {
			log.Infof("customer telphone updated: %s", telephone)
		}
	}
	_, err = collectionCustomer.UpdateByID(ctx, customer.ID, bson.M{
		"$set": bson.M{
			"email":     email,
			"updatedAt": customerUpdatedAt,
		},
	})
	if err != nil {
		log.Errorf("failed to update customer email to %s: %s", email, err)
	} else {
		log.Infof("customer email updated: %s", email)
	}

	// TODO: post processing

	return &newOrder, nil
}
