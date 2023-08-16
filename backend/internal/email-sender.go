package internal

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendEmail sends an email to the user with the verification code
func SendEmail(email string, verificationCode int) error {
	from := mail.NewEmail("Fly Spray", "no-reply@threefold.tech")
	subject := "Verify your account"
	to := mail.NewEmail("user", email)
	plainTextContent := "Verification Code"
	htmlContent := fmt.Sprintf("<strong>Verification Code: </strong>%d", verificationCode)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient("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")
	_, err := client.Send(message)
	return err
}
