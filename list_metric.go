package main

import (
	//"./config"
	//"./sendmail"
	"./emailHtml"
	"./set"
	"./structs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/aws/aws-sdk-go/service/costexplorer"

	"log"

	// "os"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"time"
)

//Report for unutilized services

var endTimePointer *time.Time
var pID *string
var periodPointer *int64
var stat *string
var reports = []structs.Report{}
var config map[string]interface{}

//var unit *string
var av string

//var st = config.Unit
func getListParam(namespace, dimensionName, dimensionValue string) *cloudwatch.ListMetricsInput {
	//fmt.Println(namespace, dimensionName, dimensionValue)
	param := &cloudwatch.ListMetricsInput{
		//   MetricName: aws.String(metric),
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
	//	unit = &st
	period := int64(config["Period"].(float64))
	periodPointer = &period
	returnData := true
	maxDataPoints := int64(config["MaxDataPoints"].(float64))
	//fmt.Println(list.Dimensions)
	metricStat := cloudwatch.MetricStat{
		Metric: &cloudwatch.Metric{ /* required */
			Dimensions: []*cloudwatch.Dimension{list.Dimensions[0]},
			MetricName: list.MetricName,
			Namespace:  list.Namespace,
		},
		Period: periodPointer, /* required */
		Stat:   stat,
		//		Unit:   unit,
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

func getCostParam() *costexplorer.GetCostAndUsageInput {
	granularity := config["Granularity"].(string)
	metric1 := "BlendedCost"
	metric2 := "UnblendedCost"
	metrics := []*string{&metric1, &metric2}
	endDate := time.Now().AddDate(0, 0, -int(config["EndTime"].(float64))).Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -int(config["StartTime"].(float64))).Format("2006-01-02")
	dateInterval := costexplorer.DateInterval{}
	dateInt := &dateInterval
	dateInt = dateInterval.SetEnd(endDate)
	dateInt = dateInterval.SetStart(startDate)
	param := costexplorer.GetCostAndUsageInput{
		Granularity: &granularity,
		Metrics:     metrics,
		TimePeriod:  dateInt,
	}
	return &param
}

func getReservationParam() *costexplorer.GetReservationUtilizationInput {
	granularity := config["Granularity"].(string)
	endDate := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -4).Format("2006-01-02")
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

func judge(res *cloudwatch.GetMetricDataOutput, threshold float64, result *cloudwatch.Metric) {
	//	fmt.Println(res)
	for i := 0; i < len(res.MetricDataResults); i++ {
		metricValue := res.MetricDataResults[i].Values
		for j := 0; j < len(metricValue); j++ {
			//fmt.Println(*metricValue[j], threshold)
			if *metricValue[j] < threshold {
				serviceName := result.Namespace
				serviceID := result.Dimensions[0].Value
				report := "Unutilized"
				utilization := *res.MetricDataResults[0].Values[0]
				timestamp := *res.MetricDataResults[0].Timestamps[0]
				r := structs.Report{*serviceName, *serviceID, report, utilization, timestamp.String()}
				reports = append(reports, r)
				//fmt.Println(reports)
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
	// fmt.Println(config)
	namespace := config["Namespace"].([]interface{})
	dimensions := config["DimensionName"].([]interface{})
	dimensionValue := config["DimensionValue"].([]interface{})
	threshold := config["Threshold"].([]interface{})
	av = config["Stat"].(string)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := cloudwatch.New(sess)
	sve := costexplorer.New(sess)
	costRes, err := sve.GetCostAndUsage(getCostParam())
	if err != nil {
		log.Fatal(err)
	}
	// reservationReport, err := sve.GetReservationUtilization(getReservationParam())
	// fmt.Println(reservationReport, err)
	// reservationreccomendation, err := sve.GetReservationPurchaseRecommendation(&costexplorer.GetReservationPurchaseRecommendationInput{Service: &serv})
	// fmt.Println(reservationreccomendation, err)
	//fmt.Println(costRes.ResultsByTime[0].Total["UnblendedCost"])
	costReport := costRes.ResultsByTime[0].Total["UnblendedCost"].Amount
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
							judge(res, threshold[j].(float64), result.Metrics[i])

						}
					}
				} else {
					judge(res, threshold[j].(float64), result.Metrics[i])

				}
			}
		}

	}
	var sr []structs.Report = set.MakeSet(reports)
	emailHtml.SendMail(sr, *costReport)
}
