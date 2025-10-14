package client

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridClient struct {
	client *sendgrid.Client
	from   *mail.Email
}

func InitializeSendGrid() *SendGridClient {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	fromEmail := os.Getenv("SENDGRID_FROM_EMAIL")

	if apiKey == "" || fromEmail == "" {
		log.Fatal("SendGrid API key or FROM email not set in environment")
	}

	client := sendgrid.NewSendClient(apiKey)
	from := mail.NewEmail("SunReady Team", fromEmail)

	return &SendGridClient{
		client: client,
		from:   from,
	}
}


func (sg *SendGridClient) SendWelcomeEmail(toEmail, name string) error {
	to := mail.NewEmail(name, toEmail)
	subject := "Welcome to SunReady!"
	plainTextContent := fmt.Sprintf("Hello %s,\n\nWelcome to SunReady! We're excited to have you on board.", name)
	htmlContent := fmt.Sprintf("<strong>Hello %s,</strong><br><br>Welcome to SunReady! We're excited to have you on board.", name)

	message := mail.NewSingleEmail(sg.from, subject, to, plainTextContent, htmlContent)
	response, err := sg.client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("failed to send email, status code: %d, body: %s", response.StatusCode, response.Body)
	}

	fmt.Printf("Welcome email sent to %s successfully\n", toEmail)
	return nil
}
