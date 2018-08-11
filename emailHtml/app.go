package emailHtml

import "../structs"

var config = Config{}

func Configuration() {
	config.Read()
}

func SendMail(report []structs.Report, costReport string) {
	Configuration()
	subject := "Daily AWS Service Report"
	destination := "sushant@exotel.in"
	r := NewRequest([]string{destination, "sushant.gupta@mountblue.io"}, subject)
	temp := make(map[string]interface{})
	temp["report"] = report
	temp["cost"] = costReport
	r.Send("emailHtml/templates/template.html", temp)
}
