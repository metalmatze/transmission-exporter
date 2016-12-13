package main

import (
	"log"
	"net/http"

	transmission "github.com/metalmatze/transmission-exporter"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	path string = "/metrics"
	addr string = ":19091"
)

func main() {
	log.Println("starting transmission-exporter")

	client := transmission.New("http://localhost:9091", nil)

	prometheus.MustRegister(NewTorrentCollector(client))

	http.Handle(path, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Node Exporter</title></head>
			<body>
			<h1>Transmission Exporter</h1>
			<p><a href="` + path + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Fatal(http.ListenAndServe(addr, nil))
}
