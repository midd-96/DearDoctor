package config

import (
	"bytes"
	"log"
	"net/smtp"
	"os"
)

type MailConfig interface {
	SendMail(to string, message string) error
}

type mailConfig struct{}

func NewMailConfig() MailConfig {
	return &mailConfig{}
}

func (c *mailConfig) SendMail(to string, message string) error {

	log.Println("Email Id to send message : ", to)
	userName := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", userName, password, smtpHost)

	headers := make(map[string]string)

	headers["Subject"] = "Dear Doctor"
	headers["From"] = userName

	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(k + ": " + v + "\r\n")
	}

	msg.WriteString("\r\n")
	msg.WriteString(message)

	// sending email
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, userName, []string{to}, msg.Bytes())

}
