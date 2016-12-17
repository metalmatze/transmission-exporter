package main

import (
	"log"
	"net/http"

	arg "github.com/alexflint/go-arg"
	"github.com/joho/godotenv"
	transmission "github.com/metalmatze/transmission-exporter"
	"github.com/prometheus/client_golang/prometheus"
)

// Config gets its content from env and passes it on to different packages
type Config struct {
	TransmissionAddr     string `arg:"env:TRANSMISSION_ADDR"`
	TransmissionPassword string `arg:"env:TRANSMISSION_PASSWORD"`
	TransmissionUsername string `arg:"env:TRANSMISSION_USERNAME"`
	WebAddr              string `arg:"env:WEB_ADDR"`
	WebPath              string `arg:"env:WEB_PATH"`
}

func main() {
	log.Println("starting transmission-exporter")

	err := godotenv.Load()
	if err != nil {
		log.Println("no .env present")
	}

	c := Config{
		WebPath:          "/metrics",
		WebAddr:          ":19091",
		TransmissionAddr: "http://localhost:9091",
	}

	arg.MustParse(&c)

	var user *transmission.User
	if c.TransmissionUsername != "" && c.TransmissionPassword != "" {
		user = &transmission.User{
			Username: c.TransmissionUsername,
			Password: c.TransmissionPassword,
		}
	}

	client := transmission.New(c.TransmissionAddr, user)

	prometheus.MustRegister(NewTorrentCollector(client))
	prometheus.MustRegister(NewSessionCollector(client))

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

func boolToString(true bool) string {
	if true {
		return "1"
	}
	return "0"
}
