package main

import (
	"./config"
	//"./sendmail"
	"./emailHtml"
	"./set"
	"./structs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/aws/aws-sdk-go/service/costexplorer"

	"fmt"
	// "os"

	"strconv"
	"time"
)

//Report for unutilized services

var endTimePointer *time.Time
var pID *string
var periodPointer *int64
var stat *string
var reports = []structs.Report{}

//var unit *string
var av = config.Stat

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
	endTime := time.Now().AddDate(0, 0, -config.EndTime)
	startTime := time.Now().AddDate(0, 0, -config.StartTime)
	startTimePointer := &startTime
	endTimePointer = &endTime
	id := "m" + strconv.Itoa(index+1)
	pID = &id
	stat = &av
	//	unit = &st
	period := int64(config.Period)
	periodPointer = &period
	returnData := true
	maxDataPoints := int64(config.MaxDataPoints)
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
	granualarity := "MONTHLY"
	metric1 := "BlendedCost"
	metric2 := "UnblendedCost"
	metric3 := "UsageQuantity"
	metric4 := "NormalizedUsageAmount"
	metric5 := "AmortizedCost"
	metric6 := "NetUnblendedCost"
	metrics := []*string{&metric1, &metric2, &metric3, &metric4, &metric5, &metric6}
	endDate := "2018-08-08"
	startDate := "2018-07-07"
	dateInterval := costexplorer.DateInterval{}
	dateInt := &dateInterval
	dateInt = dateInterval.SetEnd(endDate)
	dateInt = dateInterval.SetStart(startDate)
	param := costexplorer.GetCostAndUsageInput{
		Granularity: &granualarity,
		Metrics:     metrics,
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
	namespace := config.Namespace
	dimensions := config.DimensionName
	dimensionValue := config.DimensionValue
	threshold := config.Threshold
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := cloudwatch.New(sess)
	sve := costexplorer.New(sess)
	costRes, err := sve.GetCostAndUsage(getCostParam())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(costRes)
	for j := 0; j < len(namespace); j++ {
		result, err := svc.ListMetrics(getListParam(namespace[j], dimensions[j], dimensionValue[j]))
		if err != nil {
			fmt.Println("Error", err)
			return
		}
		for i := 0; i < len(result.Metrics); i++ {
			paramQuery := getParam(i, result.Metrics[i])
			res, err := svc.GetMetricData(&paramQuery)
			if err != nil {
				fmt.Println(i, err)
			} else {
				if res.NextToken != nil {
					for res.NextToken != nil {
						paramQuery.NextToken = res.NextToken
						res, err = svc.GetMetricData(&paramQuery)
						if err != nil {
							fmt.Println(i, err)
						} else {
							judge(res, threshold[j], result.Metrics[i])
							//								fmt.Println(res)
							// serviceName := result.Metrics[0].Namespace
							// serviceID := result.Metrics[0].Dimensions[0].Value
							// report := "Unutilized"
							// utilization := *res.MetricDataResults[0].Values[0]
							// timestamp := *res.MetricDataResults[0].Timestamps[0]
							// r := structs.Report{*serviceName, *serviceID, report, utilization, timestamp.String()}
							// reports = append(reports, r)

						}
					}
				} else {
					judge(res, threshold[j], result.Metrics[i])
					//						fmt.Println(res)

				}
			}
		}

	}
	//emailHtml.Configuration()
	var sr []structs.Report = set.MakeSet(reports)
	//fmt.Printf("%T", sr)
	/* 	data := "" */
	/* for k := 0; k < len(sr); k++ {
		data = data + "\n " + sr[k].ServiceName + " " + sr[k].ServiceID + " " + sr[k].Report + " " + sr[k].Timestamp
		//	fmt.Println(data)
	} */
	//sendmail.SendMail1(data)
	emailHtml.SendMail(sr)
}
