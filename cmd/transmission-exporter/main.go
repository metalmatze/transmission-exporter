package main

import (
	"log"
	"net/http"

	arg "github.com/alexflint/go-arg"
	transmission "github.com/metalmatze/transmission-exporter"
	"github.com/prometheus/client_golang/prometheus"
)

// Config gets its content from env and passes it on to different packages
type Config struct {
	WebPath              string `arg:"env:WEB_PATH"`
	WebAddr              string `arg:"env:WEB_ADDR"`
	TransmissionAddr     string `arg:"env:TRANSMISSION_ADDR"`
	TransmissionUsername string `arg:"env:TRANSMISSION_USERNAME"`
	TransmissionPassword string `arg:"env:TRANSMISSION_PASSWORD"`
}

func main() {
	log.Println("starting transmission-exporter")

	c := Config{
		WebPath:          "/metrics",
		WebAddr:          ":19091",
		TransmissionAddr: "http://localhost:9091",
	}

	arg.MustParse(&c)

	client := transmission.New(c.TransmissionAddr, nil)

	prometheus.MustRegister(NewTorrentCollector(client))

	http.Handle(c.WebPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Node Exporter</title></head>
			<body>
			<h1>Transmission Exporter</h1>
			<p><a href="` + c.WebPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Fatal(http.ListenAndServe(c.WebAddr, nil))
}
