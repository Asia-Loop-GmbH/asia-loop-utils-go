package servicemailjet

import (
	"context"

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
	From       Email
	To         []Email
	Subject    string
	TemplateID TemplateID
}

func Send(ctx context.Context, input SendInput, variables map[string]interface{}) error {
	receivers := make(mailjet.RecipientsV31, 0)
	for _, to := range input.To {
		receivers = append(receivers, mailjet.RecipientV31{
			Email: to.Address,
			Name:  to.Name,
		})
	}

	bcc := make(mailjet.RecipientsV31, 0)
	bcc = append(bcc, mailjet.RecipientV31{Email: ccEmail})

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
		Variables:        variables,
	}

	err := mymailjet.Send(ctx, info)

	if err != nil {
		return err
	}
	return nil
}
