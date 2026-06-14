package handlers

import (
	"fmt"
	"net/smtp"
)

func Email(to, subject, body string) error {
	from := "support.nooboj@gmail.com"
	appPassword := "cdnigpnmdrpjducm"

	message := []byte(
		"To: " + to + "\r\n" +
			"From: NoobOJ <" + from + ">\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body,
	)

	auth := smtp.PlainAuth(
		"",
		from,
		appPassword,
		"smtp.gmail.com",
	)

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		from,
		[]string{to},
		message,
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
