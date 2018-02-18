package lib

import (
	"os"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

const (
	TemplateForgotPassword = 317685
)

func MailJetConn() *mailjet.Client {
	return mailjet.NewMailjetClient(os.Getenv("MJ_APIKEY_PUBLIC"), os.Getenv("MJ_APIKEY_PRIVATE"))
}

func SendMail(mailjetClient *mailjet.Client, toEmail, toName, subject string,
	templateId interface{}, variables map[string]interface{}) error {
	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: "valentin.omnes@gmail.com",
				Name:  "Mail Delivery - Matcha",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: toEmail,
					Name:  toName,
				},
			},
			TemplateID:       templateId,
			TemplateLanguage: true,
			Subject:          subject,
			Variables:        variables,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		return err
	}
	return nil
}
