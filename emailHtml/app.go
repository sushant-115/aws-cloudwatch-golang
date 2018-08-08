package emailHtml

var config = Config{}

func configuration() {
	config.Read()
}

func SendMail() {
	configuration()
	subject := "Get latest Tech News directly to your inbox"
	destination := "mohamed.labouardy@gmail.com"
	r := NewRequest([]string{destination}, subject)
	r.Send("templates/template.html", map[string]string{"username": "Conor"})
}
