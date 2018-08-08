package config

//Namespace of the metric
var Namespace = []string{"AWS/ApplicationELB", "AWS/EC2", "AWS/EBS"}

//DimensionName for filter the data points
var DimensionName = []string{"MetricName", "MetricName", "MetricName"}

//DimensionValue for filter the data points
var DimensionValue = []string{"RequestCount", "CPUUtilization", "VolumeWriteBytes"}

//Threshold for each metric for notification
var Threshold =[]float64{100 , 10.00001, 500000}

//Stat value
var Stat = "Average"

//Unit in which the output will come
var Unit = []string{"Count", "Percent", "Bytes"}

//StartTime number of days from today
var StartTime = 1

//EndTime number of days from today (0 for today)
var EndTime = 0

//Period in seconds
var Period = 3600*24

//MaxDataPoints in result
var MaxDataPoints = 1000
