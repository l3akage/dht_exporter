package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

const version string = "0.1"

var (
	showVersion   = flag.Bool("version", false, "Print version information.")
	listenAddress = flag.String("listen-address", ":9330", "Address on which to expose metrics.")
	metricsPath   = flag.String("path", "/metrics", "Path under which to expose metrics.")
	device        = flag.Int("device", 22, "Sensor type, either 11 or 22 for DHT11/DHT22")
	gpio          = flag.Int("gpio", 4, "GPIO connected to the sensors (not pin number)")
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: dht_exporter [ ... ]\n\nParameters:")
		fmt.Println()
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	startServer()
}

func printVersion() {
	fmt.Println("dht_exporter")
	fmt.Printf("Version: %s\n", version)
}

func startServer() {
	log.Infof("Starting DHT exporter (Version: %s)\n", version)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>DHT Exporter (Version ` + version + `)</title></head>
			<body>
			<h1>DHT Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			<h2>More information:</h2>
			<p><a href="https://github.com/l3akage/dht_exporter">github.com/l3akage/dht_exporter</a></p>
			</body>
			</html>`))
	})
	http.HandleFunc(*metricsPath, handleMetricsRequest)

	log.Infof("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	reg := prometheus.NewRegistry()
	reg.MustRegister(&dhtCollector{})

	promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		ErrorLog:      log.NewErrorLogger(),
		ErrorHandling: promhttp.ContinueOnError}).ServeHTTP(w, r)
}
