package sendmail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"strings"
)

type Mail struct {
	senderId string
	toIds    []string
	subject  string
	body     string
}

type SmtpServer struct {
	host string
	port string
}

func (s *SmtpServer) ServerName() string {
	return s.host + ":" + s.port
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderId)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + MIME
	message += "\r\n" + mail.body

	return message
}

var body = `<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><meta http-equiv=\"X-UA-Compatible\" content=\"ie=edge\"><title>Document</title><style>h1{font-size: 30px;color: #fff;text-transform: uppercase;font-weight: 300;text-align: center;margin-bottom: 15px;}table{width:100%;table-layout: fixed;}.tbl-header{ background-color: rgba(255,255,255,0.3);}.tbl-content{height:300px;overflow-x:auto;margin-top: 0px;border: 1px solid rgba(255,255,255,0.3);}th{padding: 20px 15px;text-align: left;font-weight: 500;font-size: 12px;color: #fff;text-transform: uppercase;}td{padding: 15px;text-align: left;vertical-align:middle; font-weight: 300;font-size: 12px;color: #fff;border-bottom: solid 1px rgba(255,255,255,0.1);}@import url(https://fonts.googleapis.com/css?family=Roboto:400,500,300,700);body{background: -webkit-linear-gradient(left, #25c481, #25b7c4); background: linear-gradient(to right, #25c481, #25b7c4);font-family: 'Roboto', sans-serif;}section{margin: 50px;}.made-with-love {margin-top: 40px;padding: 10px;clear: left;text-align: center;font-size: 10px;font-family: arial;color: #fff;}.made-with-love i {font-style: normal;color: #F50057;font-size: 14px;position: relative;top: 2px;}.made-with-love a {  color: #fff;text-decoration: none;}.made-with-love a:hover {  text-decoration: underline;}::-webkit-scrollbar {width: 6px;}::-webkit-scrollbar-track {-webkit-box-shadow: inset 0 0 6px rgba(0,0,0,0.3);}::-webkit-scrollbar-thumb {-webkit-box-shadow: inset 0 0 6px rgba(0,0,0,0.3);}</style></head><body><section><h1>AWS Service Reports</h1><div class=\"tbl-header\"><table cellpadding=\"0\" cellspacing=\"0\" border=\"0\"><thead><tr><th>Service Name</th><th>Service ID</th><th>Report</th><th>Timestamp</th></tr></thead></table></div><div class=\"tbl-content\"><table cellpadding=\"0\" cellspacing=\"0\" border=\"0\"><tbody><tr><td>AAC</td><td>AUSTRALIAN COMPANY </td><td>$1.38</td><td>+2.01</td></tr><tr><td>AAD</td><td>AUSENCO</td><td>$2.38</td><td>-0.01</td></tr><tr><td>AAX</td><td>ADELAIDE</td><td>$3.22</td><td>+0.01</td></tr><tr><td>XXD</td><td>ADITYA BIRLA</td><td>$1.02</td><td>-1.01</td></tr><tr><td>AAC</td><td>AUSTRALIAN COMPANY </td><td>$1.38</td><td>+2.01</td></tr><tr><td>AAD</td><td>AUSENCO</td><td>$2.38</td><td>-0.01</td></tr></tbody></table></div></section></body></html>`

func parseTemplate(fileName string) string {

	t, err := template.ParseFiles(fileName)
	if err != nil {
		return ""
	}
	buffer := new(bytes.Buffer)
	var data interface{}
	if err = t.Execute(buffer, data); err != nil {
		return ""
	}
	return buffer.String()
}
func SendMail1(emailData string) {
	mail := Mail{}
	mail.senderId = "sushant78080@gmail.com"
	mail.toIds = []string{"sushant@exotel.in"}
	mail.subject = "This is the email subject"
	mail.body = parseTemplate(body)
	messageBody := mail.BuildMessage()

	smtpServer := SmtpServer{host: "smtp.gmail.com", port: "465"}

	log.Println(smtpServer.host)
	//build an auth
	auth := smtp.PlainAuth("", mail.senderId, "qazxsw321", smtpServer.host)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}

	conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		log.Panic(err)
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	// step 2: add all from and to
	if err = client.Mail(mail.senderId); err != nil {
		log.Panic(err)
	}
	for _, k := range mail.toIds {
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()

	log.Println("Mail sent successfully")

}
