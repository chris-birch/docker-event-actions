package main

import (
	"errors"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

func buildEMail(timestamp time.Time, from string, to []string, subject string, body string) string {
	var msg strings.Builder
	msg.WriteString("From: " + from + "\r\n")
	msg.WriteString("To: " + strings.Join(to, ";") + "\r\n")
	msg.WriteString("Date: " + timestamp.Format(time.RFC1123Z) + "\r\n")
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	msg.WriteString("Subject: " + subject + "\r\n")
	msg.WriteString("\r\n" + body + "\r\n")

	return msg.String()
}

func sendMail(timestamp time.Time, message string, title string, errCh chan ReporterError) {

	e := ReporterError{
		Reporter: "Mail",
	}

	from := glb_arguments.Reporter.Mail.From
	to := []string{glb_arguments.Reporter.Mail.To}
	username := glb_arguments.Reporter.Mail.User
	password := glb_arguments.Reporter.Mail.Password

	host := glb_arguments.Reporter.Mail.Host
	port := strconv.Itoa(glb_arguments.Reporter.Mail.Port)
	address := host + ":" + port

	subject := title
	body := message

	mail := buildEMail(timestamp, from, to, subject, body)

	auth := smtp.PlainAuth("", username, password, host)

	err := smtp.SendMail(address, auth, from, to, []byte(mail))
	if err != nil {
		logger.Error().Err(err).Str("reporter", "Mail").Msg("")
		e.Error = errors.New("failed to send mail")
		errCh <- e
		return
	}

}
