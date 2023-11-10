package deliverect

type OrderStatus int

const (
	// POS

	OrderStatusNew           OrderStatus = 10
	OrderStatusAccepted      OrderStatus = 20
	OrderStatusPrinted       OrderStatus = 40
	OrderStatusPreparing     OrderStatus = 50
	OrderStatusPrepared      OrderStatus = 60
	OrderStatusPickupReady   OrderStatus = 70
	OrderStatusFinalized     OrderStatus = 90
	OrderStatusAutoFinalized OrderStatus = 95
	OrderStatusCanceled      OrderStatus = 110
	OrderStatusFailed        OrderStatus = 120

	// SYSTEM

	OrderStatusParsed OrderStatus = 1
)

type OrderType int

const (
	OrderTypePickup   OrderType = 1
	OrderTypeDelivery OrderType = 2
	OrderTypeEatIn    OrderType = 3
)

type Customer struct {
	Name            *string `json:"name,omitempty"`
	CompanyName     *string `json:"companyName,omitempty"`
	PhoneNumber     *string `json:"phoneNumber,omitempty"`
	PhoneAccessCode *string `json:"phoneAccessCode,omitempty"`
	Email           *string `json:"email,omitempty"`
	Note            *string `json:"note,omitempty"`
	TIN             *string `json:"tin,omitempty"`
}

type DeliveryAddress struct {
	Street           *string `json:"street,omitempty"`
	PostalCode       *string `json:"postalCode,omitempty"`
	City             *string `json:"city,omitempty"`
	Country          *string `json:"country,omitempty"`
	Source           *string `json:"source,omitempty"`
	ExtraAddressInfo *string `json:"extraAddressInfo,omitempty"`
}

type OrderTax struct {
	Name       *string `json:"name,omitempty"`
	TaxClassID *int    `json:"taxClassId,omitempty"`
	Total      *int    `json:"total,omitempty"`
}

type OrderPaymentType int

const (
	OrderPaymentTypeCreditCardOnline OrderPaymentType = 0
	OrderPaymentTypeCash             OrderPaymentType = 1
	OrderPaymentTypeOnDelivery       OrderPaymentType = 2
	OrderPaymentTypeOnline           OrderPaymentType = 3
	OrderPaymentTypeCreditCardAtDoor OrderPaymentType = 4
	OrderPaymentTypePINAtDoor        OrderPaymentType = 5
	OrderPaymentTypeVoucherAtDoor    OrderPaymentType = 6
	OrderPaymentTypeMealVoucher      OrderPaymentType = 7
	OrderPaymentTypeBankContact      OrderPaymentType = 8
	OrderPaymentTypeOther            OrderPaymentType = 9
)

type OrderPayment struct {
	Amount int              `json:"amount"`
	Type   OrderPaymentType `json:"type"`
	Due    *int             `json:"due,omitempty"`
	Rebate *int             `json:"rebate,omitempty"`
}

type OrderItem struct {
	PLU         string       `json:"plu"`
	Name        string       `json:"name"`
	Price       int          `json:"price"`
	Quantity    int          `json:"quantity"`
	ProductType ProductType  `json:"productType"`
	Remark      *string      `json:"remark,omitempty"`
	IsInternal  *bool        `json:"isInternal,omitempty"`
	SubItems    *[]OrderItem `json:"subItems,omitempty"`
	ProductTags *[]int       `json:"productTags,omitempty"`
	IsCombo     *bool        `json:"isCombo,omitempty"`
	SortOrder   *int         `json:"sortOrder,omitempty"`
}

type DiscountType string

const (
	DiscountTypeOrderPercentOff    DiscountType = "order_percent_off"
	DiscountTypeOrderFlatOff       DiscountType = "order_flat_off"
	DiscountTypeItemBogof          DiscountType = "item_bogof"
	DiscountTypeItemFree           DiscountType = "item_free"
	DiscountTypeItemPercentOff     DiscountType = "item_percent_off"
	DiscountTypeItemFlatOff        DiscountType = "item_flat_off"
	DiscountTypeCategoryDiscount   DiscountType = "category_discount"
	DiscountTypeFreeDelivery       DiscountType = "free_delivery"
	DiscountTypeDeliveryDiscount   DiscountType = "delivery_discount"
	DiscountTypeFreeServiceFee     DiscountType = "free_service_fee"
	DiscountTypeServiceFeeDiscount DiscountType = "service_fee_discount"
	DiscountTypeStampCard          DiscountType = "stamp_card"
	DiscountTypeUnknown            DiscountType = "unknown"
)

type OrderDiscount struct {
	Type                *DiscountType `json:"type,omitempty"`
	Provider            *any          `json:"provider,omitempty"`
	Name                *string       `json:"name,omitempty"`
	ChannelDiscountCode *string       `json:"channelDiscountCode,omitempty"`
	ReferenceId         *int          `json:"referenceId,omitempty"`
	Value               *int          `json:"value,omitempty"`
	Amount              *int          `json:"amount,omitempty"`
	AmountRestaurant    *int          `json:"amountRestaurant,omitempty"`
	AmountChannel       *int          `json:"amountChannel,omitempty"`
}

type Order struct {
	Created                   string           `json:"_created"`
	Updated                   string           `json:"_updated"`
	ID                        string           `json:"_id"`
	Account                   string           `json:"account"`
	ChannelOrderID            string           `json:"channelOrderId"`
	ChannelOrderDisplayId     string           `json:"channelOrderDisplayId"`
	POSID                     *string          `json:"posId,omitempty"`
	POSReceiptID              *string          `json:"posReceiptId"`
	POSLocationID             *string          `json:"posLocationId,omitempty"`
	Location                  string           `json:"location"`
	ChannelLink               *string          `json:"channelLink,omitempty"`
	Status                    *OrderStatus     `json:"status,omitempty"`
	StatusHistory             *[]any           `json:"statusHistory,omitempty"`
	By                        *string          `json:"by,omitempty"`
	OrderType                 OrderType        `json:"orderType"`
	Channel                   int              `json:"channel"`
	POS                       *int             `json:"pos,omitempty"`
	PickupTime                *string          `json:"pickupTime,omitempty"`
	EstimatedPickupTime       *string          `json:"estimatedPickupTime,omitempty"`
	DeliveryTime              *string          `json:"deliveryTime,omitempty"`
	DeliveryIsAsap            *bool            `json:"deliveryIsAsap,omitempty"`
	Customer                  *Customer        `json:"customer,omitempty"`
	DeliveryAddress           *DeliveryAddress `json:"deliveryAddress,omitempty"`
	OrderIsAlreadyPaid        bool             `json:"orderIsAlreadyPaid"`
	TaxTotal                  *int             `json:"taxTotal,omitempty"`
	Taxes                     *[]OrderTax      `json:"taxes,omitempty"`
	TaxRemitted               *int             `json:"taxRemitted,omitempty"`
	Payment                   *OrderPayment    `json:"payment,omitempty"`
	Note                      *string          `json:"note,omitempty"`
	Items                     *[]OrderItem     `json:"items,omitempty"`
	DecimalDigits             *int             `json:"decimalDigits,omitempty"`
	NumberOfCustomers         *int             `json:"numberOfCustomers,omitempty"`
	ChannelOrderRawID         *string          `json:"channelOrderRawId,omitempty"`
	ChannelOrderHistoryRawIDs *[]string        `json:"channelOrderHistoryRawIds,omitempty"`
	ServiceCharge             *int             `json:"serviceCharge,omitempty"`
	DeliveryCost              *int             `json:"deliveryCost,omitempty"`
	BagFee                    *int             `json:"bagFee,omitempty"`
	Tip                       *int             `json:"tip,omitempty"`
	DriverTip                 *int             `json:"driverTip,omitempty"`
	DiscountTotal             *int             `json:"discountTotal,omitempty"`
	Discounts                 *[]OrderDiscount `json:"discounts,omitempty"`
	CapacityUsages            *[]string        `json:"capacityUsages,omitempty"`
	BrandId                   *string          `json:"brandId,omitempty"`
	TestOrder                 *bool            `json:"testOrder,omitempty"`
	Timezone                  *string          `json:"timezone,omitempty"`
	Date                      *int             `json:"date,omitempty"` // Date In format YMD
	Tags                      *[]string        `json:"tags,omitempty"`
}
