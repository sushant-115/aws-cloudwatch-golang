package emailHtml

import "aws-cloudwatch-golang/structs"

//ReportCount contains all counts of all unutilized service and their name
type ReportCount struct {
	Name  string
	Count int
}

//SendMail will send the mail to recipients mentioned in configuration.json
func SendMail(report []structs.Report, additionalReport *structs.AdditionalReport) {
	request := NewRequest(additionalReport.MailRecipients, additionalReport.MailSubject)
	mapToSendToHTMLMail := make(map[string]interface{})
	var reportCountArr []ReportCount
	reportCount := make(map[string]int)
	for i := 0; i < len(report); i++ {
		reportCount[report[i].ServiceName]++
	}
	for k := range reportCount {
		rept := ReportCount{k, reportCount[k]}
		reportCountArr = append(reportCountArr, rept)
	}
	mapToSendToHTMLMail["report"] = report
	mapToSendToHTMLMail["cost"] = additionalReport.DailyCostReport
	mapToSendToHTMLMail["reportCount"] = reportCountArr
	mapToSendToHTMLMail["unusedHour"] = additionalReport.RIUnusedHours
	mapToSendToHTMLMail["utilization"] = additionalReport.RIUtilizationPercentage
	mapToSendToHTMLMail["monthlyCost"] = additionalReport.MonthlyCostReport
	mapToSendToHTMLMail["startDate"] = additionalReport.StartDate
	mapToSendToHTMLMail["endDate"] = additionalReport.EndDate
	request.Send("emailHtml/templates/template.html", mapToSendToHTMLMail)
}
