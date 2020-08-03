package main

import (
	"log"
	"time"

	"github.com/metalmatze/transmission-exporter"
	"github.com/prometheus/client_golang/prometheus"
)

// SessionStatsCollector exposes SessionStats as metrics
type SessionStatsCollector struct {
	clients []*transmission.Client

	DownloadSpeed  *prometheus.Desc
	UploadSpeed    *prometheus.Desc
	TorrentsTotal  *prometheus.Desc
	TorrentsActive *prometheus.Desc
	TorrentsPaused *prometheus.Desc

	Downloaded   *prometheus.Desc
	Uploaded     *prometheus.Desc
	FilesAdded   *prometheus.Desc
	ActiveTime   *prometheus.Desc
	SessionCount *prometheus.Desc
}

// NewSessionStatsCollector takes a transmission.Client and returns a SessionStatsCollector
func NewSessionStatsCollector(clients []*transmission.Client) *SessionStatsCollector {
	const collectorNamespace = "session_stats_"

	return &SessionStatsCollector{
		clients: clients,

		DownloadSpeed: prometheus.NewDesc(
			namespace+collectorNamespace+"download_speed_bytes",
			"Current download speed in bytes",
			[]string{"client_name"},
			nil,
		),
		UploadSpeed: prometheus.NewDesc(
			namespace+collectorNamespace+"upload_speed_bytes",
			"Current download speed in bytes",
			[]string{"client_name"},
			nil,
		),
		TorrentsTotal: prometheus.NewDesc(
			namespace+collectorNamespace+"torrents_total",
			"The total number of torrents",
			[]string{"client_name"},
			nil,
		),
		TorrentsActive: prometheus.NewDesc(
			namespace+collectorNamespace+"torrents_active",
			"The number of active torrents",
			[]string{"client_name"},
			nil,
		),
		TorrentsPaused: prometheus.NewDesc(
			namespace+collectorNamespace+"torrents_paused",
			"The number of paused torrents",
			[]string{"client_name"},
			nil,
		),

		Downloaded: prometheus.NewDesc(
			namespace+collectorNamespace+"downloaded_bytes",
			"The number of downloaded bytes",
			[]string{"type", "client_name"},
			nil,
		),
		Uploaded: prometheus.NewDesc(
			namespace+collectorNamespace+"uploaded_bytes",
			"The number of uploaded bytes",
			[]string{"type", "client_name"},
			nil,
		),
		FilesAdded: prometheus.NewDesc(
			namespace+collectorNamespace+"files_added",
			"The number of files added",
			[]string{"type", "client_name"},
			nil,
		),
		ActiveTime: prometheus.NewDesc(
			namespace+collectorNamespace+"active",
			"The time transmission is active since",
			[]string{"type", "client_name"},
			nil,
		),
		SessionCount: prometheus.NewDesc(
			namespace+collectorNamespace+"sessions",
			"Count of the times transmission started",
			[]string{"type", "client_name"},
			nil,
		),
	}
}

// Describe implements the prometheus.Collector interface
func (sc *SessionStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- sc.DownloadSpeed
	ch <- sc.UploadSpeed
	ch <- sc.TorrentsTotal
	ch <- sc.TorrentsActive
	ch <- sc.TorrentsPaused
}

// Collect implements the prometheus.Collector interface
func (sc *SessionStatsCollector) Collect(ch chan<- prometheus.Metric) {
	for _, client := range sc.clients {
		stats, err := client.GetSessionStats()
		if err != nil {
			log.Printf("failed to get session stats, client: %v, error: %v", client.Name, err)
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			sc.DownloadSpeed,
			prometheus.GaugeValue,
			float64(stats.DownloadSpeed),
			client.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.UploadSpeed,
			prometheus.GaugeValue,
			float64(stats.UploadSpeed),
			client.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.TorrentsTotal,
			prometheus.GaugeValue,
			float64(stats.TorrentCount),
			client.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.TorrentsActive,
			prometheus.GaugeValue,
			float64(stats.ActiveTorrentCount),
			client.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.TorrentsPaused,
			prometheus.GaugeValue,
			float64(stats.PausedTorrentCount),
			client.Name,
		)

		types := []string{"current", "cumulative"}
		for _, t := range types {
			var stateStats transmission.SessionStateStats
			if t == types[0] {
				stateStats = stats.CurrentStats
			} else {
				stateStats = stats.CumulativeStats
			}

			ch <- prometheus.MustNewConstMetric(
				sc.Downloaded,
				prometheus.GaugeValue,
				float64(stateStats.DownloadedBytes),
				t, client.Name,
			)
			ch <- prometheus.MustNewConstMetric(
				sc.Uploaded,
				prometheus.GaugeValue,
				float64(stateStats.UploadedBytes),
				t, client.Name,
			)
			ch <- prometheus.MustNewConstMetric(
				sc.FilesAdded,
				prometheus.GaugeValue,
				float64(stateStats.FilesAdded),
				t, client.Name,
			)

			dur := time.Duration(stateStats.SecondsActive) * time.Second
			timestamp := time.Now().Add(-1 * dur).Unix()

			ch <- prometheus.MustNewConstMetric(
				sc.ActiveTime,
				prometheus.GaugeValue,
				float64(timestamp),
				t, client.Name,
			)
			ch <- prometheus.MustNewConstMetric(
				sc.SessionCount,
				prometheus.GaugeValue,
				float64(stateStats.SessionCount),
				t, client.Name,
			)
		}
	}
}
