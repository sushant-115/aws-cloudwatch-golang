package main

import (
	"./config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"fmt"
	// "os"

	"strconv"
	"time"
)

var endTimePointer *time.Time
var pID *string
var periodPointer *int64
var stat *string

//var unit *string
var av = config.Stat

//var st = config.Unit
func getListParam(namespace, dimensionName, dimensionValue string) *cloudwatch.ListMetricsInput {
	fmt.Println(namespace, dimensionName, dimensionValue)
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

func main() {
	namespace := config.Namespace
	dimensions := config.DimensionName
	dimensionValue := config.DimensionValue

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := cloudwatch.New(sess)
	for j := 0; j < len(namespace); j++ {
		result, err := svc.ListMetrics(getListParam(namespace[j], dimensions[j], dimensionValue[j]))
		if err != nil {
			fmt.Println("Error", err)
			return
		}
		fmt.Println(result)
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
							fmt.Println(res)
						}
					}
				} else {
					fmt.Println(res)
				}
			}
		}
	}
}
