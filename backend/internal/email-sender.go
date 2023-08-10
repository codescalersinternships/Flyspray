package internal

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendEmail sends an email to the user with the verification code
func SendEmail(email, verificationCode string) error {
	from := mail.NewEmail("Fly Spray", os.Getenv("EMAIL"))
	subject := "Verify your account"
	to := mail.NewEmail("user", email)
	plainTextContent := "Verification Code"
	htmlContent := fmt.Sprintf("<strong>Verification Code: </strong>%s", verificationCode)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	return err
}
