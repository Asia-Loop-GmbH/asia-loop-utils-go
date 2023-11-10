package deliverect

type ProductPriceLevels struct {
	TA *int `json:"TA,omitempty"`
	DL *int `json:"DL,omitempty"`
	UE *int `json:"UE,omitempty"`
}

type NameTranslations struct {
	EN *string `json:"en,omitempty"`
	ES *string `json:"es,omitempty"`
	FR *string `json:"fr,omitempty"`
	NL *string `json:"nl,omitempty"`
	AR *string `json:"ar,omitempty"`
	EL *string `json:"el,omitempty"`
}

type NutritionalInfo struct {
	Fat           *int `json:"fat,omitempty"`
	Sugar         *int `json:"sugar,omitempty"`
	SaturatedFat  *int `json:"saturatedFat,omitempty"`
	Carbohydrates *int `json:"carbohydrates,omitempty"`
	Protein       *int `json:"protein,omitempty"`
	Salt          *int `json:"salt,omitempty"`
}

type SupplementalInfo struct {
	InstructionsForUse *string   `json:"instructionsForUse,omitempty"`
	Ingredients        *[]string `json:"ingredients,omitempty"`
	Additives          *[]string `json:"additives,omitempty"`
	Prepackaged        *bool     `json:"prepackaged,omitempty"`
	Deposit            *int      `json:"deposit,omitempty"`
}

type Category struct {
	Name          string `json:"name"`
	POSCategoryID string `json:"posCategoryId"`
}

type ProductType int

const (
	ProductTypeProduct       ProductType = 1
	ProductTypeModifier      ProductType = 2
	ProductTypeModifierGroup ProductType = 3
	ProductTypeBundle        ProductType = 4
)

type Product struct {
	Name               string              `json:"name"`
	PLU                string              `json:"plu"`
	Price              int                 `json:"price"`
	Description        *string             `json:"description,omitempty"`
	ProductType        ProductType         `json:"productType"`
	IsVariant          *bool               `json:"isVariant,omitempty"`
	IsCombo            *bool               `json:"isCombo,omitempty"`
	DeliveryTax        int                 `json:"deliveryTax"`
	TakeawayTax        int                 `json:"takeawayTax"`
	EatInTax           int                 `json:"eatInTax"`
	PriceLevels        *ProductPriceLevels `json:"priceLevels,omitempty"`
	Overloads          *[]any              `json:"overloads,omitempty"`
	POSProductID       *string             `json:"posProductId,omitempty"`
	POSCategoryIDs     []string            `json:"posCategoryIds"`
	ImageURL           *string             `json:"imageUrl,omitempty"`
	NameTranslations   *NameTranslations   `json:"nameTranslations,omitempty"`
	ProductTags        *[]int              `json:"productTags,omitempty"`
	MultiMax           int                 `json:"multiMax"`
	Min                int                 `json:"min"`
	Max                int                 `json:"max"`
	DefaultQuantity    *int                `json:"defaultQuantity,omitempty"`
	NutritionalInfo    *NutritionalInfo    `json:"nutritionalInfo,omitempty"`
	SupplementalInfo   *SupplementalInfo   `json:"supplementalInfo,omitempty"`
	BottleDepositPrice *int                `json:"bottleDepositPrice,omitempty"`
	Visible            *bool               `json:"visible,omitempty"`
	SubProducts        []string            `json:"subProducts,omitempty"`
}

type PriceLevels struct {
	Name  *string `json:"name,omitempty"`
	POSID *string `json:"posId,omitempty"`
}
