package normalizer

import (
	"context"
	"regexp"
	"strings"

	"github.com/Masterminds/goutils"
	"github.com/nyaruka/phonenumbers"

	"github.com/nam-truong-le/lambda-utils-go/v2/pkg/logger"
)

func Email(ctx context.Context, email string) string {
	log := logger.FromContext(ctx)
	result := strings.ToLower(email)
	log.Infof("normalized email: '%s' -> '%s'", email, result)
	return result
}

func PhoneNumber(ctx context.Context, phone string) string {
	log := logger.FromContext(ctx)
	parsed, err := phonenumbers.Parse(phone, "DE")
	if err != nil {
		log.Warnf("could not normalize phone '%s', return original", phone)
		return phone
	}
	result := phonenumbers.Format(parsed, phonenumbers.INTERNATIONAL)
	log.Infof("normalized phone: '%s' -> '%s'", phone, result)
	return result
}

func Name(ctx context.Context, name string) string {
	log := logger.FromContext(ctx)
	result := normalizeWhitespace(ctx, goutils.CapitalizeFully(name, ' ', '-'))
	log.Infof("normalized name: '%s' -> '%s'", name, result)
	return result
}

func normalizeWhitespace(ctx context.Context, text string) string {
	space, err := regexp.Compile(`\s+`)
	if err != nil {
		return text
	}
	duplicateRemoved := space.ReplaceAllString(text, " ")
	return strings.Trim(duplicateRemoved, " ")
}
