package emailHtml

import "../structs"

var config = Config{}

func Configuration() {
	config.Read()
}

func SendMail(report []structs.Report) {
	Configuration()
	subject := "Get latest Tech News directly to your inbox"
	destination := "sushant@exotel.in"
	r := NewRequest([]string{destination}, subject)
	r.Send("emailHtml/templates/template.html", map[string][]structs.Report{"report": report})
}
