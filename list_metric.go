package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"fmt"
	// "os"
	"./config"
	"strconv"
	"time"
)

var t *time.Time
var pID *string
var period *int64
var stat *string
var unit *string
var av = config.Stat
var st = config.Unit

func getParam(index int, list *cloudwatch.Metric) cloudwatch.GetMetricDataInput {
	tim := time.Now()
	startTime := time.Now().AddDate(0, 0, -2)
	stime := &startTime
	t = &tim
	id := "m" + strconv.Itoa(index+1)
	pID = &id
	stat = &av
	unit = &st
	prd := int64(3000)
	period = &prd
	returnData := true
	maxDataPoints := int64(1000)
	metricStat := cloudwatch.MetricStat{
		Metric: &cloudwatch.Metric{ /* required */
			Dimensions: []*cloudwatch.Dimension{list.Dimensions[0]},
			MetricName: list.MetricName,
			Namespace:  list.Namespace,
		},
		Period: period, /* required */
		Stat:   stat,
		Unit:   unit,
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
	//    metric := os.Args[1]
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
		}
		fmt.Println(res)
	}
}
