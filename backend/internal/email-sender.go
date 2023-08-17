package internal

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendEmail sends an email to the user with the verification code
func SendEmail(apiKey,apiEmail, email string, verificationCode int) error {
	from := mail.NewEmail("Fly Spray", apiEmail)
	subject := "Verify your account"
	to := mail.NewEmail("user", email)
	plainTextContent := "Verification Code"
	htmlContent := fmt.Sprintf("<strong>Verification Code: </strong>%d", verificationCode)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(apiKey)
	_, err := client.Send(message)
	return err
}
