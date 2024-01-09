package mailservice

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendEmail(emailBody string) (bool, error) {
	host := os.Getenv("MAIL_SERVICE_HOST")
	user := os.Getenv("MAIL_SERVICE_USER")
	pass := os.Getenv("MAIL_SERVICE_PASS")
	port, _ := strconv.Atoi(os.Getenv("MAIL_SERVICE_PORT"))

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", user)
	mailer.SetHeader("To", user)
	mailer.SetHeader("Subject", "The Tip of The Day")
	mailer.SetBody("text/html", emailBody)
	dialer := gomail.NewDialer(host, port, user, pass)
	err := dialer.DialAndSend(mailer)

	if err != nil {

		return false, err
	}

	return true, nil
}
