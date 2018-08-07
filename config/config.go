package config

//Namespace of the metric
const Namespace string = "AWS/Billing"

//DimensionName for filter the data points
const DimensionName string = "Currency"

//DimensionValue for filter the data points
const DimensionValue string = ""

//Stat value
const Stat string = "Average"

//Unit in which the output will come
const Unit string = "Percent"

//StartTime number of days from today
const StartTime int = 7

//EndTime number of days from today (0 for today)
const EndTime int = 0

//Period in seconds
const Period int = 300

//MaxDataPoints in result
const MaxDataPoints int = 1000
