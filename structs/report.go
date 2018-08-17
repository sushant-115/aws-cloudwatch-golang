package structs

//Report struct for the generated reports
type Report struct {
	ServiceName  string
	ServiceID    string
	Report       string
	Utiliization string
	Timestamp    string
}

//AdditionalReport contains additional information for generating report
type AdditionalReport struct {
	DailyCostReport         string
	MonthlyCostReport       string
	RIUnusedHours           string
	RIUtilizationPercentage string
	MailRecipients          []string
	MailSubject             string
	StartDate               string
	EndDate                 string
}
