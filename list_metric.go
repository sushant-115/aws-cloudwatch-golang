package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"aws-cloudwatch-golang/emailHtml"
	"aws-cloudwatch-golang/set"
	"aws-cloudwatch-golang/structs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

var endTimePointer *time.Time
var pID *string
var periodPointer *int64
var stat *string
var reports = []structs.Report{}
var config map[string]interface{}

var av string

func getListParam(namespace, dimensionName, dimensionValue string) *cloudwatch.ListMetricsInput {
	param := &cloudwatch.ListMetricsInput{

		Namespace: aws.String(namespace),
		Dimensions: []*cloudwatch.DimensionFilter{
			&cloudwatch.DimensionFilter{
				Name:  aws.String(dimensionName),
				Value: aws.String(dimensionValue),
			},
		},
	}
	return param
}

func getParam(index int, list *cloudwatch.Metric) cloudwatch.GetMetricDataInput {
	if len(list.Dimensions) < 1 {
		return cloudwatch.GetMetricDataInput{}
	}
	endTime := time.Now().AddDate(0, 0, -int(config["EndTime"].(float64)))
	startTime := time.Now().AddDate(0, 0, -int(config["StartTime"].(float64)))
	startTimePointer := &startTime
	endTimePointer = &endTime
	id := "m" + strconv.Itoa(index+1)
	pID = &id
	stat = &av

	period := int64(config["Period"].(float64))
	periodPointer = &period
	returnData := true
	maxDataPoints := int64(config["MaxDataPoints"].(float64))

	metricStat := cloudwatch.MetricStat{
		Metric: &cloudwatch.Metric{ /* required */
			Dimensions: []*cloudwatch.Dimension{list.Dimensions[0]},
			MetricName: list.MetricName,
			Namespace:  list.Namespace,
		},
		Period: periodPointer, /* required */
		Stat:   stat,
	}
	metricQuery := cloudwatch.MetricDataQuery{
		Id:         pID, /* required */
		MetricStat: &metricStat,
		ReturnData: &returnData,
	}
	param := cloudwatch.GetMetricDataInput{
		EndTime: endTimePointer, /* required */
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			&metricQuery,
		},
		StartTime:     startTimePointer,
		MaxDatapoints: &maxDataPoints,
	}
	return param
}

func getCostParam() (*costexplorer.GetCostAndUsageInput, *costexplorer.GetCostAndUsageInput) {
	granularity := config["Granularity"].(string)
	metric1 := "BlendedCost"
	metric2 := "UnblendedCost"
	t := time.Now()
	metrics := []*string{&metric1, &metric2}
	endDate := t.AddDate(0, 0, -int(config["EndTime"].(float64))).Format("2006-01-02")
	startDate := t.AddDate(0, 0, -int(config["StartTime"].(float64))).Format("2006-01-02")
	dateInterval := costexplorer.DateInterval{}
	dateInt := &dateInterval
	dateInt = dateInterval.SetEnd(endDate)
	dateInt = dateInterval.SetStart(startDate)
	param := costexplorer.GetCostAndUsageInput{
		Granularity: &granularity,
		Metrics:     metrics,
		TimePeriod:  dateInt,
	}
	dateIntervalMonth := costexplorer.DateInterval{}
	dateIntMonth := &dateIntervalMonth
	dateIntMonth = dateIntervalMonth.SetEnd(endDate)
	monthFirstDay := t.AddDate(0, 0, -t.Day()+1).Format("2006-01-02")
	dateIntMonth = dateIntervalMonth.SetStart(monthFirstDay)
	param2 := costexplorer.GetCostAndUsageInput{
		Granularity: &granularity,
		Metrics:     metrics,
		TimePeriod:  dateIntMonth,
	}
	return &param, &param2
}

func getReservationParam() *costexplorer.GetReservationUtilizationInput {
	granularity := config["Granularity"].(string)
	endDate := time.Now().AddDate(0, 0, -3-(int(config["EndTime"].(float64)))).Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -3-(int(config["StartTime"].(float64)))).Format("2006-01-02")
	dateInterval := costexplorer.DateInterval{}
	dateInt := &dateInterval
	dateInt = dateInterval.SetEnd(endDate)
	dateInt = dateInterval.SetStart(startDate)
	param := costexplorer.GetReservationUtilizationInput{
		Granularity: &granularity,
		TimePeriod:  dateInt,
	}
	return &param
}

func judge(res *cloudwatch.GetMetricDataOutput, threshold float64, result *cloudwatch.Metric, utilSuffix string) {

	for i := 0; i < len(res.MetricDataResults); i++ {
		metricValue := res.MetricDataResults[i].Values
		for j := 0; j < len(metricValue); j++ {

			if *metricValue[j] < threshold {

				serviceName := result.Namespace
				serviceID := result.Dimensions[0].Value
				report := "Unutilized"
				utilization := strconv.FormatFloat(*res.MetricDataResults[0].Values[0], 'g', -1, 32) + " " + utilSuffix
				timestamp := *res.MetricDataResults[0].Timestamps[0]
				r := structs.Report{*serviceName, *serviceID, report, utilization, timestamp.String()}
				reports = append(reports, r)

			}
		}
	}

}

func main() {
	bytesJson, er := ioutil.ReadFile("configuration.json")
	if er != nil {
		log.Fatal(er)
	}
	er = json.Unmarshal(bytesJson, &config)
	if er != nil {
		log.Fatal(er)
	}

	namespace := config["Namespace"].([]interface{})
	dimensions := config["DimensionName"].([]interface{})
	dimensionValue := config["DimensionValue"].([]interface{})
	threshold := config["Threshold"].([]interface{})
	mailRecipients := config["MailRecipients"].([]interface{})
	subject := config["Subject"].(string)
	suffixs := config["Suffix"].([]interface{})
	startDate := time.Now().AddDate(0, 0, -int(config["StartTime"].(float64))).Format("2006-01-02")
	endDate := time.Now().AddDate(0, 0, -int(config["EndTime"].(float64))).Format("2006-01-02")
	var mailRecipientsStr []string
	for i := 0; i < len(mailRecipients); i++ {
		mailRecipientsStr = append(mailRecipientsStr, mailRecipients[i].(string))
	}
	av = config["Stat"].(string)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := cloudwatch.New(sess)
	sve := costexplorer.New(sess)
	dailyParam, monthlyParam := getCostParam()
	costRes, err := sve.GetCostAndUsage(dailyParam)
	if err != nil {
		log.Fatal(err)
	}
	monthlyCost, err := sve.GetCostAndUsage(monthlyParam)

	reservationReport, err := sve.GetReservationUtilization(getReservationParam())

	costReport := costRes.ResultsByTime[0].Total["UnblendedCost"].Amount
	costReportMonth := monthlyCost.ResultsByTime[0].Total["UnblendedCost"].Amount
	for j := 0; j < len(namespace); j++ {
		result, err := svc.ListMetrics(getListParam(namespace[j].(string), dimensions[j].(string), dimensionValue[j].(string)))
		if err != nil {
			log.Fatal(err)
			return
		}
		for i := 0; i < len(result.Metrics); i++ {
			paramQuery := getParam(i, result.Metrics[i])
			res, err := svc.GetMetricData(&paramQuery)
			if err != nil {
				log.Fatal(err)
			} else {
				if res.NextToken != nil {
					for res.NextToken != nil {
						paramQuery.NextToken = res.NextToken
						res, err = svc.GetMetricData(&paramQuery)
						if err != nil {
							log.Fatal(err)
						} else {
							judge(res, threshold[j].(float64), result.Metrics[i], suffixs[j].(string))

						}
					}
				} else {
					judge(res, threshold[j].(float64), result.Metrics[i], suffixs[j].(string))

				}
			}
		}

	}
	var tableReport []structs.Report = set.MakeSet(reports)
	additionalReport := structs.AdditionalReport{*costReport, *costReportMonth, *reservationReport.Total.UnusedHours, *reservationReport.Total.UtilizationPercentage, mailRecipientsStr, subject, startDate, endDate}
	emailHtml.SendMail(tableReport, &additionalReport)
}
