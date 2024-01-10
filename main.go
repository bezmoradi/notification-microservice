package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bezmoradi/notification-microservice/consumer"
	"github.com/bezmoradi/notification-microservice/helpers"
	"github.com/bezmoradi/notification-microservice/mailservice"
)

func init() {
	helpers.LoadEnvironmentVariables()
}

func main() {
	reader := consumer.Reader()

	defer reader.Close()

	for {
		message, err := reader.ReadMessage(context.Background())

		if err != nil {
			fmt.Println("Error reading message", err)
		}

		sendEmail(string(message.Value))
	}

}

func sendEmail(messageBody string) {
	emailBody, emailBodyIsValid := helpers.JsonValidator(messageBody)

	if !emailBodyIsValid {
		fmt.Println("Email body is invalid!")

		return
	}

	html := "<!DOCTYPE html><html><head>"
	styles, _ := os.ReadFile("styles.txt")
	html += string(styles)
	html += "</head><body>"
	html += emailBody
	html += "</body></html>"

	_, err := mailservice.SendEmail(html)

	if err != nil {
		fmt.Println("Error sending email", err)

		return
	}

	fmt.Println("Email sent successfully")
}
