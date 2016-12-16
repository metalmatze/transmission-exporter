package main

import (
	"log"
	"strconv"

	transmission "github.com/metalmatze/transmission-exporter"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace string = "transmission_"
)

// TorrentCollector has a transmission.Client to create torrent metrics
type TorrentCollector struct {
	client *transmission.Client

	Status   *prometheus.Desc
	Added    *prometheus.Desc
	Files    *prometheus.Desc
	Finished *prometheus.Desc
	Done     *prometheus.Desc
	Ratio    *prometheus.Desc
	Download *prometheus.Desc
	Upload   *prometheus.Desc

	// TrackerStats
	Downloads *prometheus.Desc
	Leechers  *prometheus.Desc
	Seeders   *prometheus.Desc
}

// NewTorrentCollector creates a new torrent collector with the transmission.Client
func NewTorrentCollector(client *transmission.Client) *TorrentCollector {
	return &TorrentCollector{
		client: client,

		Status: prometheus.NewDesc(
			namespace+"torrent_status",
			"Status of a torrent",
			[]string{"id", "name"},
			nil,
		),
		Added: prometheus.NewDesc(
			namespace+"torrent_added",
			"The unixtime time a torrent was added",
			[]string{"id", "name"},
			nil,
		),
		Files: prometheus.NewDesc(
			namespace+"torrent_files_total",
			"The unixtime time a torrent was added",
			[]string{"id", "name"},
			nil,
		),
		Finished: prometheus.NewDesc(
			namespace+"torrent_finished",
			"Indicates if a torrent is finished (1) or not (0)",
			[]string{"id", "name"},
			nil,
		),
		Done: prometheus.NewDesc(
			namespace+"torrent_done",
			"The percent of a torrent being done",
			[]string{"id", "name"},
			nil,
		),
		Ratio: prometheus.NewDesc(
			namespace+"torrent_ratio",
			"The upload ratio of a torrent",
			[]string{"id", "name"},
			nil,
		),
		Download: prometheus.NewDesc(
			namespace+"torrent_download_bytes",
			"The current download rate of a torrent in bytes",
			[]string{"id", "name"},
			nil,
		),
		Upload: prometheus.NewDesc(
			namespace+"torrent_upload_bytes",
			"The current upload rate of a torrent in bytes",
			[]string{"id", "name"},
			nil,
		),

		// TrackerStats
		Downloads: prometheus.NewDesc(
			namespace+"torrent_downloads_total",
			"",
			[]string{"id", "name", "tracker"},
			nil,
		),
		Leechers: prometheus.NewDesc(
			namespace+"torrent_leechers",
			"",
			[]string{"id", "name", "tracker"},
			nil,
		),
		Seeders: prometheus.NewDesc(
			namespace+"torrent_seeders",
			"",
			[]string{"id", "name", "tracker"},
			nil,
		),
	}
}

// Describe implements the prometheus.Collector interface
func (tc *TorrentCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- tc.Ratio
}

// Collect implements the prometheus.Collector interface
func (tc *TorrentCollector) Collect(ch chan<- prometheus.Metric) {
	torrents, err := tc.client.GetTorrents()
	if err != nil {
		log.Printf("failed to get torrents: %v", err)
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
			tc.Added,
			prometheus.GaugeValue,
			float64(t.Added),
			id, t.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			tc.Files,
			prometheus.GaugeValue,
			float64(len(t.Files)),
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

		for _, tracker := range t.TrackerStats {
			ch <- prometheus.MustNewConstMetric(
				tc.Downloads,
				prometheus.GaugeValue,
				float64(tracker.DownloadCount),
				id, t.Name, tracker.Host,
			)

			ch <- prometheus.MustNewConstMetric(
				tc.Leechers,
				prometheus.GaugeValue,
				float64(tracker.LeecherCount),
				id, t.Name, tracker.Host,
			)

			ch <- prometheus.MustNewConstMetric(
				tc.Seeders,
				prometheus.GaugeValue,
				float64(tracker.SeederCount),
				id, t.Name, tracker.Host,
			)
		}
	}
}
