package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"os"
	"strconv"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/mail/commons"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/users"
	"github.com/matcornic/hermes/v2"

	gomail "gopkg.in/mail.v2"
)

var (

	//	Initialize commons variables
	FROM     string
	PASSWORD string
	HOST     string
	PORT     int
)

type Service interface {
	Invite(context.ServiceContext, *commons.InvitationOptions) *errors.Error
	SendKey(context.ServiceContext, *commons.SendKeyOptions) *errors.Error
}

type DefaultMailService struct{}

func (*DefaultMailService) Invite(ctx context.ServiceContext, options *commons.InvitationOptions) *errors.Error {

	//	Initialize commons variables
	FROM = os.Getenv("SMTP_USERNAME")
	PASSWORD = os.Getenv("SMTP_PASSWORD")
	HOST = os.Getenv("SMTP_HOST")
	PORT, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", FROM)

	// Set E-Mail receivers
	m.SetHeader("To", options.ReceiverEmail)

	// Set E-Mail subject
	m.SetHeader("Subject", "Invitation to Join Organisation")

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	Fetch the user who sent the invite.
	sender, err := users.Get(ctx, client, options.SenderID)
	if err != nil {
		return err
	}

	//	Fetch the user who is receiver the invite.
	receiver, err := users.GetByEmail(ctx, client, options.ReceiverEmail)
	if err != nil {
		return err
	}

	//	Fetch the organisation for which the invite has been sent.
	organisation, err := organisations.Get(ctx, client, options.OrgID)
	if err != nil {
		return err
	}

	email := hermes.Email{
		Body: hermes.Body{
			Greeting: "Hey",
			Name:     receiver.Name,
			Intros: []string{
				fmt.Sprintf("You have been invited by %s to join the organisation: %s", sender.Name, organisation.Name),
			},
			Actions: []hermes.Action{
				{
					Instructions: "Accepting this invite requires you to have an envsecrets account. If you do not already have one, please create one from app.envsecrets.com.",
					Button: hermes.Button{
						Color: "#222", // Optional action button color
						Text:  "Accept Invite",
						Link:  fmt.Sprintf("%s/v1/invites/accept?key=%s&id=%s", os.Getenv("API"), options.Key, options.ID),
					},
				},
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	body, er := commons.Hermes.GenerateHTML(email)
	if er != nil {
		return errors.New(er, "failed to generate html for email", errors.ErrorTypeEmailFailed, errors.ErrorSourceHermes)
	}

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", body)

	// Settings for SMTP server
	d := gomail.NewDialer(HOST, PORT, FROM, PASSWORD)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		return errors.New(err, "failed to send invitation email", errors.ErrorTypeEmailFailed, errors.ErrorSourceMailer)
	}

	return nil
}

func (*DefaultMailService) SendKey(ctx context.ServiceContext, options *commons.SendKeyOptions) *errors.Error {

	//	Initialize commons variables
	FROM = os.Getenv("SMTP_USERNAME")
	PASSWORD = os.Getenv("SMTP_PASSWORD")
	HOST = os.Getenv("SMTP_HOST")
	PORT, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))

	//	Initialize Hasura client with admin privileges
	client := clients.NewGQLClient(&clients.GQLConfig{
		Type: clients.HasuraClientType,
		Headers: []clients.Header{
			clients.XHasuraAdminSecretHeader,
		},
	})

	//	Fetch the user who sent the invite.
	user, err := users.Get(ctx, client, options.UserID)
	if err != nil {
		return err
	}

	email := hermes.Email{
		Body: hermes.Body{
			Greeting: "Hey",
			Name:     user.Name,
			Intros: []string{
				"Congratulations on creating your new organisation: " + options.OrgName,
				"We've generated a fresh encryption key for you. All the secrets you save/generate in this organisation will be encrypted with this new key.",
				"You are strongly advised to keep this key somewhere safe.",
			},
			Outros: []string{
				"Kindly find the key file attached with this email.",
				"You can also manually export this key from your organisation's settings page.",
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	body, er := commons.Hermes.GenerateHTML(email)
	if er != nil {
		return errors.New(er, "failed to generate html for email", errors.ErrorTypeEmailFailed, errors.ErrorSourceHermes)
	}

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", FROM)

	// Set E-Mail receivers
	m.SetHeader("To", user.Email)

	// Set E-Mail subject
	m.SetHeader("Subject", "Encryption Key For New Organisation")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", body)

	//	Attach the key.
	b := bytes.NewBufferString(options.Key)
	m.AttachReader("key.txt", b)

	// Settings for SMTP server
	d := gomail.NewDialer(HOST, PORT, FROM, PASSWORD)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		return errors.New(err, "failed to send invitation email", errors.ErrorTypeEmailFailed, errors.ErrorSourceMailer)
	}

	return nil
}
