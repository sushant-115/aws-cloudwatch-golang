package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudwatch"

    "fmt"
    "os"
)

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

    fmt.Println("Metrics", result.Metrics[0].Namespace)
}
