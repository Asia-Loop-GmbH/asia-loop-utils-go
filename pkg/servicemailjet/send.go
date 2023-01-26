package servicemailjet

import (
	"context"
	"encoding/json"

	"github.com/mailjet/mailjet-apiv3-go/v4"

	mymailjet "github.com/nam-truong-le/lambda-utils-go/v3/pkg/mailjet"
)

type TemplateID int

const (
	TemplateIDOrder TemplateID = 3889572
	ccEmail                    = "order@asialoop.de"
)

type Email struct {
	Address string
	Name    string
}

type SendInput struct {
	From        Email
	To          []Email
	Subject     string
	TemplateID  TemplateID
	Attachments *mailjet.AttachmentsV31
}

func Send[T any](ctx context.Context, input SendInput, variables T) error {
	receivers := make(mailjet.RecipientsV31, 0)
	for _, to := range input.To {
		receivers = append(receivers, mailjet.RecipientV31{
			Email: to.Address,
			Name:  to.Name,
		})
	}

	bcc := make(mailjet.RecipientsV31, 0)
	bcc = append(bcc, mailjet.RecipientV31{Email: ccEmail})

	mapVariables := make(map[string]interface{})
	v, err := json.Marshal(variables)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(v, &mapVariables); err != nil {
		return err
	}

	info := mailjet.InfoMessagesV31{
		From: &mailjet.RecipientV31{
			Email: input.From.Address,
			Name:  input.From.Name,
		},
		To:               &receivers,
		Bcc:              &bcc,
		Subject:          input.Subject,
		TemplateID:       int(input.TemplateID),
		TemplateLanguage: true,
		Variables:        mapVariables,
		Attachments:      input.Attachments,
	}

	if err := mymailjet.Send(ctx, info); err != nil {
		return err
	}

	return nil
}
