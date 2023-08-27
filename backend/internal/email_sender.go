package internal

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendEmail sends an email to the user
func SendEmail(apiKey, sender, receiver, subject, body string) error {
	from := mail.NewEmail("Fly Spray", sender)
	to := mail.NewEmail("user", receiver)
	message := mail.NewSingleEmail(from, subject, to, "", body)

	client := sendgrid.NewSendClient(apiKey)
	_, err := client.Send(message)
	return err
}

// VerifyMailContent generates subject and body for mail verification
func VerifyMailContent(verificationCode int) (string, string) {
	subject := "Verify your account"
	body := fmt.Sprintf("<strong>Verification Code: </strong>%d", verificationCode)
	return subject, body
}

// ForgetPasswordMailContent generates subject and body for forget password mail
func ForgetPasswordMailContent(verificationCode int) (string, string) {
	subject := "Forget Password Request"
	body := fmt.Sprintf("<strong>Verification Code: </strong>%d", verificationCode)
	return subject, body
}
