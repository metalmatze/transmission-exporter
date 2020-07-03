package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	arg "github.com/alexflint/go-arg"
	"github.com/joho/godotenv"
	transmission "github.com/metalmatze/transmission-exporter"
	"github.com/metalmatze/transmission-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config gets its content from env and passes it on to different packages
type Config struct {
	TransmissionAddr     string `arg:"env:TRANSMISSION_ADDR"`
	TransmissionPassword string `arg:"env:TRANSMISSION_PASSWORD"`
	TransmissionUsername string `arg:"env:TRANSMISSION_USERNAME"`
	ClientName           string `arg:"env:CLIENT_NAME"`
	WebAddr              string `arg:"env:WEB_ADDR"`
	WebPath              string `arg:"env:WEB_PATH"`
	ConfigFile           string `arg:"env:CONFIG_FILE"`
}

func main() {
	log.Println("starting transmission-exporter")

	c, err := loadConfig()

	if err != nil {
		log.Fatalf("Could not load config: %v", err)
		os.Exit(3)
	}

	var clients []*transmission.Client

	for _, clientConf := range c.Clients {
		client := transmission.New(clientConf)

		clients = append(clients, client)
	}

	prometheus.MustRegister(NewTorrentCollector(clients))
	prometheus.MustRegister(NewSessionCollector(clients))
	prometheus.MustRegister(NewSessionStatsCollector(clients))

	http.Handle(c.WebPath, promhttp.Handler())

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

func loadConfig() (*config.Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env present")
	}

	c := Config{
		WebPath:          "/metrics",
		WebAddr:          ":19091",
		TransmissionAddr: "http://localhost:9091",
	}

	arg.Parse(&c)

	// load config from file
	if c.ConfigFile != "" {
		b, err := ioutil.ReadFile(c.ConfigFile)
		if err != nil {
			return nil, err
		}

		return config.Load(bytes.NewReader(b))
	}

	// load config from flags or env
	if c.ClientName == "" {
		c.ClientName = c.TransmissionAddr
	}
	return &config.Config{
		WebAddr: c.WebAddr,
		WebPath: c.WebPath,
		Clients: []*config.Client{
			&config.Client{
				ClientName:           c.ClientName,
				TransmissionAddr:     c.TransmissionAddr,
				TransmissionUsername: c.TransmissionUsername,
				TransmissionPassword: c.TransmissionPassword,
			},
		},
	}, nil
}
