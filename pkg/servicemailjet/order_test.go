package servicemailjet_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/asia-loop-gmbh/asia-loop-utils-go/v6/pkg/servicemailjet"
	commoncontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
)

func TestSendOrder(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.WithValue(context.TODO(), commoncontext.FieldStage, "dev")
	err := servicemailjet.Send[servicemailjet.SendOrderVariables](ctx, servicemailjet.SendInput{
		From: servicemailjet.Email{
			Address: "noreply@asialoop.de",
			Name:    "Asia Loop GmbH",
		},
		To: []servicemailjet.Email{{
			Address: "lenamtruong@gmail.com",
			Name:    "Nam Truong Le",
		}},
		Subject:    "Neue Bestellung #1234",
		TemplateID: servicemailjet.TemplateIDOrder,
	}, servicemailjet.SendOrderVariables{
		FirstName:     "Nam Truong",
		Title:         "vielen Dank für Deine Bestellung",
		Content:       "Wir haben Deine Bestellung empfangen. Bei Fragen oder Anmerkungen zu Deiner Bestellung möchten wir Dich bitten, das Restaurant telefonisch zu kontaktieren. Für weniger dringende Fragen kannst Du mit unserem Kundenservice Kontakt aufnehmen.",
		ActionText:    "Verfolge Deine Bestellung",
		ActionLink:    "https://google.de",
		ActionEnabled: "true",
	})

	assert.NoError(t, err)
}
