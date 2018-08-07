package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudwatch"

    "fmt"
   // "os"
    "strconv"
    "time"
)
var t *time.Time
var pId *string
var period *int64
var stat *string
var unit *string
var av string = "Average"
var st string = "Seconds"
func getParam (index int,list *cloudwatch.Metric) cloudwatch.GetMetricDataInput{
	tim := time.Now()
        t = &tim
        id := "m"+strconv.Itoa(index+1)
        pId = &id
        stat = &av
        unit = &st
        prd := int64(3000)
        period = &prd
        param:= cloudwatch.GetMetricDataInput{
	EndTime: t , /* required */
  MetricDataQueries: []*cloudwatch.MetricDataQuery  {
      Id: pId, /* required */
     MetricStat:*cloudwatch.MetricStat{
       Metric:&cloudwatch.Metric { /* required */
          Dimensions:*cloudwatch.Dimension{list.Dimensions},
          MetricName:list.MetricName,
          Namespace: list.Namespace,
        },
        Period:period, /* required */
        Stat: stat,
        Unit :unit,
      },
      ReturnData: true,
    },
  StartTime:t,
  MaxDatapoints: 1000,        
        }
  return param
}

func main() {
//    metric := os.Args[1]
    namespace := "AWS/EBS"
    dimensions := "VolumeId"
    
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

     svc := cloudwatch.New(sess)

     result, err := svc.ListMetrics(&cloudwatch.ListMetricsInput{
     //   MetricName: aws.String(metric),
        Namespace:  aws.String(namespace),
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

    for i:=0;i<len(result.Metrics);i++ {
      res ,err := svc.GetMetricsData(getParam(i,result.Metrics[i]))
      if err!=nil {
      fmt.Println(i , "err")
      }
	fmt.Println(res)
    }
}
