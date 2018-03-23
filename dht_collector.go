package main

import (
	"fmt"
	"os"

	dht "github.com/d2r2/go-dht"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix = "dht_"

var (
	upDesc       *prometheus.Desc
	tempDesc     *prometheus.Desc
	humidityDesc *prometheus.Desc
)

func init() {
	upDesc = prometheus.NewDesc(prefix+"up", "Scrape was successful", nil, nil)
	tempDesc = prometheus.NewDesc(prefix+"temp", "Air temperature (in degrees C)", nil, nil)
	humidityDesc = prometheus.NewDesc(prefix+"humidity_percent", "Humidity", nil, nil)
}

type dhtCollector struct {
}

func (c dhtCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upDesc
	ch <- tempDesc
	ch <- humidityDesc
}

func (c dhtCollector) Collect(ch chan<- prometheus.Metric) {
	sensor := dht.DHT22
	if *device == 11 {
		sensor = dht.DHT11
	}

	temperature, humidity, _, err := dht.ReadDHTxxWithRetry(sensor, *gpio, false, 10)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error getting sensor data", err)
		ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, 1)
		ch <- prometheus.MustNewConstMetric(tempDesc, prometheus.GaugeValue, float64(temperature))
		ch <- prometheus.MustNewConstMetric(humidityDesc, prometheus.GaugeValue, float64(humidity))
	}
}
