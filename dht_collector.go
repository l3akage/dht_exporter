package main

import (
	"fmt"
	"os"
	"strconv"

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
	l := []string{"location"}
	upDesc = prometheus.NewDesc(prefix+"up", "Scrape was successful", l, nil)
	tempDesc = prometheus.NewDesc(prefix+"temp", "Air temperature (in degrees C)", l, nil)
	humidityDesc = prometheus.NewDesc(prefix+"humidity_percent", "Humidity", l, nil)
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
	for gpiostr, name := range list.Names {
		gpio, err := strconv.Atoi(gpiostr)
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid gpiod", err)
			ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, 0, name)
			continue
		}
		temperature, humidity, _, err := dht.ReadDHTxxWithRetry(sensor, gpio, false, 10)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error getting sensor data", err)
			ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, 0, name)
			continue
		}
		ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, 1, name)
		ch <- prometheus.MustNewConstMetric(tempDesc, prometheus.GaugeValue, float64(int(temperature*100))/100, name)
		ch <- prometheus.MustNewConstMetric(humidityDesc, prometheus.GaugeValue, float64(int(humidity*100))/100, name)
	}
}
