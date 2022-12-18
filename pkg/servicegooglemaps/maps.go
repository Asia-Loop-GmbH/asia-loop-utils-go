package servicegooglemaps

import (
	"context"
	"fmt"

	"googlemaps.github.io/maps"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
)

type ResolveAddressResult struct {
	StreetNumber     string
	Street           string
	City             string
	Postcode         string
	State            string
	FormattedAddress string
}

func ResolveAddress(ctx context.Context, address string) (*ResolveAddressResult, error) {
	log := logger.FromContext(ctx)
	log.Infof("resolve address: %s", address)
	apiKey, err := ssm.GetParameter(ctx, "/google/maps/key", true)
	if err != nil {
		return nil, err
	}
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	result, err := client.Geocode(context.Background(), &maps.GeocodingRequest{
		Address:  address,
		Region:   "DE",
		Language: "DE",
	})
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("could not resolve address: %s", address)
	}

	if len(result) > 1 {
		log.Warnf("multiple addresses found, the first item will be taken")
	}

	return &ResolveAddressResult{
		StreetNumber:     getAddressComponent(result[0], "street_number"),
		Street:           getAddressComponent(result[0], "route"),
		City:             getAddressComponent(result[0], "locality"),
		Postcode:         getAddressComponent(result[0], "postal_code"),
		State:            getAddressComponent(result[0], "administrative_area_level_1"),
		FormattedAddress: result[0].FormattedAddress,
	}, nil
}

func getAddressComponent(result maps.GeocodingResult, componentName string) string {
	for _, component := range result.AddressComponents {
		for _, componentType := range component.Types {
			if componentType == componentName {
				return component.LongName
			}
		}
	}
	return ""
}
