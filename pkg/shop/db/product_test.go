package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v8/pkg/shop/db"
)

func TestIsAvailableInStore(t *testing.T) {
	tests := []struct {
		name            string
		product         db.Product
		storeKey        string
		selectedOptions map[string]string
		options         []db.ProductOption
		expectedIsAvail bool
	}{
		{
			name: "gift card product not available",
			product: db.Product{
				IsGiftCard: true,
			},
			storeKey:        "store123",
			selectedOptions: nil,
			options:         nil,
			expectedIsAvail: false,
		},
		{
			name: "disabled product not available",
			product: db.Product{
				DisabledIn: []string{"store123"},
			},
			storeKey:        "store123",
			selectedOptions: nil,
			options:         nil,
			expectedIsAvail: false,
		},
		{
			name: "out of stock product not available",
			product: db.Product{
				OutOfStockIn: []string{"store123"},
			},
			storeKey:        "store123",
			selectedOptions: nil,
			options:         nil,
			expectedIsAvail: false,
		},
		{
			name:     "selected option disabled not available",
			product:  db.Product{},
			storeKey: "store123",
			selectedOptions: map[string]string{
				"Color": "Red",
			},
			options: []db.ProductOption{
				{
					Name: "Color",
					Values: []db.ProductOptionValue{
						{
							Name:       "Red",
							DisabledIn: []string{"store123"},
						},
					},
				},
			},
			expectedIsAvail: false,
		},
		{
			name:     "selected option out of stock not available",
			product:  db.Product{},
			storeKey: "store123",
			selectedOptions: map[string]string{
				"Color": "Red",
			},
			options: []db.ProductOption{
				{
					Name: "Color",
					Values: []db.ProductOptionValue{
						{
							Name:         "Red",
							OutOfStockIn: []string{"store123"},
						},
					},
				},
			},
			expectedIsAvail: false,
		},
		{
			name:     "product available",
			product:  db.Product{},
			storeKey: "store123",
			selectedOptions: map[string]string{
				"Color": "Red",
				"Size":  "M",
			},
			options: []db.ProductOption{
				{
					Name: "Color",
					Values: []db.ProductOptionValue{
						{Name: "Red"},
					},
				},
				{
					Name: "Size",
					Values: []db.ProductOptionValue{
						{Name: "M"},
					},
				},
			},
			expectedIsAvail: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.product.IsAvailableInStore(tc.storeKey, tc.selectedOptions, tc.options)
			assert.Equal(t, tc.expectedIsAvail, result)
		})
	}
}
