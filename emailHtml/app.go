package emailHtml

var config = Config{}

func Configuration() {
	config.Read()
}

func SendMail() {
	Configuration()
	subject := "Get latest Tech News directly to your inbox"
	destination := "sushant@exotel.in"
	r := NewRequest([]string{destination}, subject)
	r.Send("emailHtml/templates/template.html", map[string]string{"username": "Sushant"})
}
