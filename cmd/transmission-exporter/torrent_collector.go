package main

import (
	"log"
	"strconv"

	transmission "github.com/metalmatze/transmission-exporter"
	"github.com/prometheus/client_golang/prometheus"
)

type TorrentCollector struct {
	client   *transmission.Client
	Status   *prometheus.Desc
	Finished *prometheus.Desc
	Done     *prometheus.Desc
	Added    *prometheus.Desc
	Ratio    *prometheus.Desc
	Download *prometheus.Desc
	Upload   *prometheus.Desc
}

func NewTorrentCollector(client *transmission.Client) *TorrentCollector {
	return &TorrentCollector{
		client: client,

		Status: prometheus.NewDesc(
			"transmission_torrent_status",
			"Status of a torrent",
			[]string{"id", "name"},
			nil,
		),
		Finished: prometheus.NewDesc(
			"transmission_torrent_finished",
			"Indicates if a torrent is finished (1) or not (0)",
			[]string{"id", "name"},
			nil,
		),
		Done: prometheus.NewDesc(
			"transmission_torrent_done",
			"The percent of a torrent being done",
			[]string{"id", "name"},
			nil,
		),
		Added: prometheus.NewDesc(
			"transmission_torrent_added",
			"The unixtime time a torrent was added",
			[]string{"id", "name"},
			nil,
		),
		Ratio: prometheus.NewDesc(
			"transmission_torrent_ratio",
			"The upload ratio of a torrent",
			[]string{"id", "name"},
			nil,
		),
		Download: prometheus.NewDesc(
			"transmission_torrent_download_bytes",
			"The current download rate of a torrent in bytes",
			[]string{"id", "name"},
			nil,
		),
		Upload: prometheus.NewDesc(
			"transmission_torrent_upload_bytes",
			"The current upload rate of a torrent in bytes",
			[]string{"id", "name"},
			nil,
		),
	}
}

func (tc *TorrentCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- tc.Ratio
}

func (tc *TorrentCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("fetching torrents...")
	torrents, err := tc.client.GetTorrents()
	if err != nil {
		return
	}

	for _, t := range torrents {
		var finished float64

		id := strconv.Itoa(t.ID)

		if t.IsFinished {
			finished = 1
		}

		ch <- prometheus.MustNewConstMetric(
			tc.Status,
			prometheus.GaugeValue,
			float64(t.Status),
			id, t.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			tc.Finished,
			prometheus.GaugeValue,
			finished,
			id, t.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			tc.Done,
			prometheus.GaugeValue,
			t.PercentDone,
			id, t.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			tc.Added,
			prometheus.GaugeValue,
			float64(t.Added),
			id, t.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			tc.Ratio,
			prometheus.GaugeValue,
			t.UploadRatio,
			id, t.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			tc.Download,
			prometheus.GaugeValue,
			float64(t.RateDownload),
			id, t.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			tc.Upload,
			prometheus.GaugeValue,
			float64(t.RateUpload),
			id, t.Name,
		)
	}
}
