package channels

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/ritchieridanko/apotekly/services/notification/internal/models"
	"github.com/ritchieridanko/apotekly/services/notification/internal/templates"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/infra/mailer"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"gopkg.in/gomail.v2"
)

type EmailChannel interface {
	SendWelcome(ctx context.Context, data *models.WelcomeEmail) (err *ce.Error)
}

type emailChannel struct {
	clientAddr string
	sender     string
	logoURL    string
	mailer     *mailer.Mailer
	template   *template.Template
}

func NewEmailChannel(clientAddr, sender, logoURL string, m *mailer.Mailer) (EmailChannel, error) {
	tmpl, err := template.ParseFS(templates.Email, "emails/*.html.tmpl", "emails/partials/*.html.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize email channel: %w", err)
	}
	return &emailChannel{
		clientAddr: clientAddr,
		sender:     sender,
		logoURL:    logoURL,
		mailer:     m,
		template:   tmpl,
	}, nil
}

func (c *emailChannel) SendWelcome(ctx context.Context, data *models.WelcomeEmail) *ce.Error {
	// URL Generation
	var url string
	if data.IsEmailVerified {
		url = c.clientAddr
	} else {
		tokenizedURL, err := utils.GenerateTokenizedURL(
			c.clientAddr,
			"/auth/verify-email",
			data.VerificationToken,
		)
		if err != nil {
			return ce.NewError(ce.CodeURLGenerationFailed, ce.MsgInternalServer, err)
		}
		url = tokenizedURL
	}

	// Template Building
	body, err := c.buildTemplate(
		"welcome",
		map[string]any{
			"Subject":         "Welcome Aboard!",
			"Recipient":       data.Recipient,
			"Title":           "Welcome to Apotekly",
			"LogoURL":         c.logoURL,
			"URL":             url,
			"Year":            time.Now().UTC().Year(),
			"IsEmailVerified": data.IsEmailVerified,
		},
	)
	if err != nil {
		return err
	}

	// Message Composition
	msg := c.composeMessage([]string{data.Recipient}, "Welcome Aboard!", body.String())

	// Email Delivery
	return c.send(msg)
}

func (c *emailChannel) buildTemplate(template string, data map[string]any) (bytes.Buffer, *ce.Error) {
	var buf bytes.Buffer
	if err := c.template.ExecuteTemplate(&buf, template, data); err != nil {
		return bytes.Buffer{}, ce.NewError(
			ce.CodeEmailTemplatingFailed,
			ce.MsgInternalServer,
			err,
			logger.NewField("email_template", template),
		)
	}
	return buf, nil
}

func (c *emailChannel) composeMessage(recipients []string, subject, body string) *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetHeader("From", c.sender)
	msg.SetHeader("To", recipients...)
	msg.SetHeader("Subject", utils.ToMIMEBase64(subject))
	msg.SetBody("text/plain", "Please view this email in an HTML-compatible client!")
	msg.AddAlternative("text/html", body)
	return msg
}

func (c *emailChannel) send(msg *gomail.Message) *ce.Error {
	if err := c.mailer.Send(msg); err != nil {
		return ce.NewError(ce.CodeEmailDeliveryFailed, ce.MsgInternalServer, err)
	}
	return nil
}
