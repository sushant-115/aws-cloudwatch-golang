package emailHtml

import "../structs"

var config = Config{}

type ReportCount struct {
	Name  string
	Count int
}

func Configuration() {
	config.Read()
}

func SendMail(report []structs.Report, costReport string) {
	Configuration()
	subject := "Daily AWS Service Report"
	destination := "sushant@exotel.in"
	r := NewRequest([]string{destination, "sushant.gupta@mountblue.io"}, subject)
	temp := make(map[string]interface{})
	var reportCountArr []ReportCount
	rc := make(map[string]int)
	for i := 0; i < len(report); i++ {
		rc[report[i].ServiceName]++
	}
	for k := range rc {
		r := ReportCount{k, rc[k]}
		reportCountArr = append(reportCountArr, r)
	}
	temp["report"] = report
	temp["cost"] = costReport
	temp["reportCount"] = reportCountArr
	r.Send("emailHtml/templates/template.html", temp)
}
