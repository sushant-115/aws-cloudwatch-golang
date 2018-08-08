package config

//Namespace of the metric
const Namespace string = "AWS/ApplicationELB"

//DimensionName for filter the data points
const DimensionName string = "LoadBalancer"

//DimensionValue for filter the data points
const DimensionValue string = ""

//Stat value
const Stat string = "Average"

//Unit in which the output will come
const Unit string = "Count"

//StartTime number of days from today
const StartTime int = 1

//EndTime number of days from today (0 for today)
const EndTime int = 0

//Period in seconds
const Period int = 300

//MaxDataPoints in result
const MaxDataPoints int = 1000
