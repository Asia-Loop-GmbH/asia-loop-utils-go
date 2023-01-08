package servicepusher

import (
	"context"
	"fmt"

	"github.com/nam-truong-le/lambda-utils-go/v2/pkg/logger"
	"github.com/nam-truong-le/lambda-utils-go/v2/pkg/pusher"
)

const (
	eventOrderCreatedName           = "order-created"
	eventOrderPOSPaymentStartedName = "order-payment-started"
	eventOrderPOSPaymentPaidName    = "order-paid"
	eventOrderDeliveredName         = "order-delivered"
	eventGroupFinalizedName         = "group-finalized"
	eventGroupDeliveredName         = "group-delivered"
)

type EventOrderCreatedData struct {
	StoreID string `json:"storeId"`
}

func PublishOrderCreated(ctx context.Context, data *EventOrderCreatedData) error {
	log := logger.FromContext(ctx)
	c, err := pusher.NewClient(ctx)
	if err != nil {
		return err
	}
	if err := c.Trigger(storeChannel(data.StoreID), eventOrderCreatedName, data); err != nil {
		log.Errorf("failed to publish event [%s]: %v", eventOrderCreatedName, data)
		return err
	}
	log.Infof("event published [%s]: %v", eventOrderCreatedName, data)
	return nil
}

type EventOrderDeliveredData struct {
	StoreID string `json:"storeId"`
}

func PublishOrderDelivered(ctx context.Context, data *EventOrderDeliveredData) error {
	log := logger.FromContext(ctx)
	c, err := pusher.NewClient(ctx)
	if err != nil {
		return err
	}
	if err := c.Trigger(storeChannel(data.StoreID), eventOrderDeliveredName, data); err != nil {
		log.Errorf("failed to publish event [%s]: %v", eventOrderDeliveredName, data)
		return err
	}
	log.Infof("event published [%s]: %v", eventOrderDeliveredName, data)
	return nil
}

type EventOrderPOSPaymentStartedData struct {
	StoreID string `json:"storeId"`
	OrderID string `json:"orderId"`
}

func PublishOrderPOSPaymentStarted(ctx context.Context, data *EventOrderPOSPaymentStartedData) error {
	log := logger.FromContext(ctx)
	c, err := pusher.NewClient(ctx)
	if err != nil {
		return err
	}
	if err := c.Trigger(storeChannel(data.StoreID), eventOrderPOSPaymentStartedName, data); err != nil {
		log.Errorf("failed to publish event [%s]: %v", eventOrderPOSPaymentStartedName, data)
		return err
	}
	log.Infof("event published [%s]: %v", eventOrderPOSPaymentStartedName, data)
	return nil
}

type EventOrderPOSPaymentPaidData struct {
	StoreID string `json:"storeId"`
	OrderID string `json:"orderId"`
}

func PublishOrderPOSPaymentPaid(ctx context.Context, data *EventOrderPOSPaymentPaidData) error {
	log := logger.FromContext(ctx)
	c, err := pusher.NewClient(ctx)
	if err != nil {
		return err
	}
	if err := c.Trigger(storeChannel(data.StoreID), eventOrderPOSPaymentPaidName, data); err != nil {
		log.Errorf("failed to publish event [%s]: %v", eventOrderPOSPaymentPaidName, data)
		return err
	}
	log.Infof("event published [%s]: %v", eventOrderPOSPaymentPaidName, data)
	return nil
}

type EventGroupFinalizedData struct {
	StoreID string `json:"storeId"`
}

func PublishGroupFinalized(ctx context.Context, data *EventGroupFinalizedData) error {
	log := logger.FromContext(ctx)
	c, err := pusher.NewClient(ctx)
	if err != nil {
		return err
	}
	if err := c.Trigger(storeChannel(data.StoreID), eventGroupFinalizedName, data); err != nil {
		log.Errorf("failed to publish event [%s]: %v", eventGroupFinalizedName, data)
		return err
	}
	log.Infof("event published [%s]: %v", eventGroupFinalizedName, data)
	return nil
}

type EventGroupDeliveredData struct {
	StoreID string `json:"storeId"`
}

func PublishGroupDelivered(ctx context.Context, data *EventGroupDeliveredData) error {
	log := logger.FromContext(ctx)
	c, err := pusher.NewClient(ctx)
	if err != nil {
		return err
	}
	if err := c.Trigger(storeChannel(data.StoreID), eventGroupDeliveredName, data); err != nil {
		log.Errorf("failed to publish event [%s]: %v", eventGroupDeliveredName, data)
		return err
	}
	log.Infof("event published [%s]: %v", eventGroupDeliveredName, data)
	return nil
}

func storeChannel(storeID string) string {
	return fmt.Sprintf("restaurant-%s", storeID)
}
