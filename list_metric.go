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

var t *time.Time
var pID *string
var period *int64
var stat *string
//var unit *string
var av = config.Stat
//var st = config.Unit

func getParam(index int, list *cloudwatch.Metric) cloudwatch.GetMetricDataInput {
	tim := time.Now().AddDate(0, 0, -config.EndTime)
	startTime := time.Now().AddDate(0, 0, -config.StartTime)
	stime := &startTime
	t = &tim
	id := "m" + strconv.Itoa(index+1)
	pID = &id
	stat = &av
//	unit = &st
	prd := int64(config.Period)
	period = &prd
	returnData := true
	maxDataPoints := int64(config.MaxDataPoints)
	metricStat := cloudwatch.MetricStat{
		Metric: &cloudwatch.Metric{ /* required */
			Dimensions: []*cloudwatch.Dimension{list.Dimensions[0]},
			MetricName: list.MetricName,
			Namespace:  list.Namespace,
		},
		Period: period, /* required */
		Stat:   stat,
//		Unit:   unit,
	}
	metricQuery := cloudwatch.MetricDataQuery{
		Id:         pID, /* required */
		MetricStat: &metricStat,
		ReturnData: &returnData,
	}
	param := cloudwatch.GetMetricDataInput{
		EndTime: t, /* required */
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			&metricQuery,
		},
		StartTime:     stime,
		MaxDatapoints: &maxDataPoints,
	}
	return param
}

func main() {
	namespace := config.Namespace
	dimensions := config.DimensionName

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := cloudwatch.New(sess)

	result, err := svc.ListMetrics(&cloudwatch.ListMetricsInput{
		//   MetricName: aws.String(metric),
		Namespace: aws.String(namespace),
		Dimensions: []*cloudwatch.DimensionFilter{
			&cloudwatch.DimensionFilter{
				Name: aws.String(dimensions),
			},
		},
	})
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
						fmt.Println(res)
					}
				}
			} else {
				fmt.Println(res)
			}
		}
	}
}
