package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/mail/commons"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/users"
	userCommons "github.com/envsecrets/envsecrets/internal/users/commons"
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
	Invite(context.ServiceContext, *commons.InvitationOptions) error
	SendKey(context.ServiceContext, *commons.SendKeyOptions) error
	SendWelcomeEmail(context.ServiceContext, *userCommons.User) error
}

type DefaultMailService struct{}

func (*DefaultMailService) Invite(ctx context.ServiceContext, options *commons.InvitationOptions) error {

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

	//	Fetch the organisation for which the invite has been sent.
	organisation, err := organisations.Get(ctx, client, options.OrgID)
	if err != nil {
		return err
	}

	email := hermes.Email{
		Body: hermes.Body{
			Greeting: "Hey",
			Intros: []string{
				"You have been invited to join a new organisation on envsecrets.",
			},
			Dictionary: []hermes.Entry{
				{Key: "Organisation", Value: organisation.Name},
				{Key: "Invited By", Value: sender.DisplayName},
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
	body, err := commons.Hermes.GenerateHTML(email)
	if err != nil {
		return err
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
		return err
	}

	return nil
}

func (*DefaultMailService) SendKey(ctx context.ServiceContext, options *commons.SendKeyOptions) error {

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
			Name:     user.DisplayName,
			Intros: []string{
				"Congratulations on creating your new organisation: " + options.OrgName,
				"We've generated a fresh encryption key for you. All the secrets you save/generate in this organisation will be encrypted with this new key.",
				"You are strongly advised to keep this key somewhere safe.",
			},
			Outros: []string{
				"Kindly find the key file attached with this email.",
			},
			Signature: "Mrinal Wahal",
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	body, err := commons.Hermes.GenerateHTML(email)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "Mrinal Wahal")

	// Set E-Mail receivers
	m.SetHeader("To", user.Email)

	// Set E-Mail subject
	m.SetHeader("Subject", "Encryption Key For Your Organisation")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", body)

	//	Attach the key.
	b := bytes.NewBufferString(options.Key)
	m.AttachReader("key.txt", b)

	// Settings for SMTP server
	d := gomail.NewDialer(HOST, PORT, FROM, PASSWORD)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	isDevEnvironment, err := strconv.ParseBool(os.Getenv("DEV"))
	if err != nil {
		return err
	}

	d.TLSConfig = &tls.Config{InsecureSkipVerify: isDevEnvironment}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (*DefaultMailService) SendWelcomeEmail(ctx context.ServiceContext, user *userCommons.User) error {

	//	Initialize commons variables
	FROM = os.Getenv("SMTP_USERNAME")
	PASSWORD = os.Getenv("SMTP_PASSWORD")
	HOST = os.Getenv("SMTP_HOST")
	PORT, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))

	email := hermes.Email{
		Body: hermes.Body{
			Greeting: "Hey",
			Name:     strings.Split(user.Name, " ")[0],
			Intros: []string{
				"I’m excited to have you here.",
				"Let me quickly get you up and running with your envsecrets account.",
				"In 5 minutes, you’ll be set with your secrets and ready to synchronise them continuously and safely.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Use the quickstart guide below to begin.",
					Button: hermes.Button{
						Color: "#222", // Optional action button color
						Text:  "Quickstart Guide",
						Link:  "https://docs.envsecrets.com/platform/quickstart",
					},
				},
				{
					Instructions: "If you have any queries, I will be right here 👇🏼",
					Button: hermes.Button{
						Color: "#222", // Optional action button color
						Text:  "Community Server",
						Link:  os.Getenv("DISCORD"),
					},
				},
			},
			Signature: "Mrinal Wahal",
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	body, err := commons.Hermes.GenerateHTML(email)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", FROM)

	// Set E-Mail receivers
	m.SetHeader("To", user.Email)

	// Set E-Mail subject
	m.SetHeader("Subject", "Welcome to envsecrets!")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", body)

	// Settings for SMTP server
	d := gomail.NewDialer(HOST, PORT, FROM, PASSWORD)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	isDevEnvironment, err := strconv.ParseBool(os.Getenv("DEV"))
	if err != nil {
		return err
	}

	d.TLSConfig = &tls.Config{InsecureSkipVerify: isDevEnvironment}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
