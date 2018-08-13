package emailHtml

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
)

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func NewRequest(to []string, subject string) *Request {
	return &Request{
		to:      to,
		subject: subject,
	}
}

func (r *Request) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *Request) sendMail() bool {
	if len(r.to) == 0 {
		fmt.Println("Minimum one recipient is neccesary")
		os.Exit(1)
	}
	var recipients string
	for i := 0; i < len(r.to); i++ {
		recipients = recipients + r.to[i]
		if i+1 != len(r.to) {
			recipients = recipients + ";"
		}
	}
	body := "To: " + recipients + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTP := fmt.Sprintf("%s:%d", config.Server, config.Port)
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", config.Email, config.Password, config.Server), config.Email, r.to, []byte(body)); err != nil {
		return false
	}
	return true
}

func (r *Request) Send(templateName string, items interface{}) {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		log.Fatal(err)
	}
	if ok := r.sendMail(); ok {
		log.Printf("Email has been sent to %s\n", r.to)
	} else {
		log.Printf("Failed to send the email to %s\n", r.to)
	}
}
