package main

import (
	"flag"
	"github.com/givanov/rpi_export/pkg/mbox"
	"github.com/givanov/rpi_export/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	flagAddr  = flag.String("addr", "", "Listen on address")
	flagDebug = flag.Bool("debug", false, "Print debug messages")
)

func main() {
	flag.Parse()
	mbox.Debug = *flagDebug

	collector := metrics.NewRaspberryPiMboxCollector()
	prometheus.MustRegister(collector)

	if *flagAddr != "" {
		http.Handle("/metrics", promhttp.Handler())
		log.Printf("Listening on %s", *flagAddr)
		http.ListenAndServe(*flagAddr, nil)
		return
	}
}
