package logger

import (
	"github.com/smira/go-statsd"
)

var statsdClient *statsd.Client

func Init() {
	statsdClient = statsd.NewClient("localhost:8125",
		statsd.MaxPacketSize(1400),
		statsd.MetricPrefix("webapp."))
}

func GetMetricsClient() *statsd.Client {
	statsdClient.Incr("endpoint.", 1)

	return statsdClient
}
