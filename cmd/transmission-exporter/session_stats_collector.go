package main

import (
	"log"
	"time"

	"github.com/metalmatze/transmission-exporter"
	"github.com/prometheus/client_golang/prometheus"
)

// SessionStatsCollector exposes SessionStats as metrics
type SessionStatsCollector struct {
	client *transmission.Client

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
func NewSessionStatsCollector(client *transmission.Client) *SessionStatsCollector {
	const collectorNamespace = "session_stats_"

	return &SessionStatsCollector{
		client: client,

		DownloadSpeed: prometheus.NewDesc(
			namespace+collectorNamespace+"download_speed_bytes",
			"Current download speed in bytes",
			nil,
			nil,
		),
		UploadSpeed: prometheus.NewDesc(
			namespace+collectorNamespace+"upload_speed_bytes",
			"Current download speed in bytes",
			nil,
			nil,
		),
		TorrentsTotal: prometheus.NewDesc(
			namespace+collectorNamespace+"torrents_total",
			"The total number of torrents",
			nil,
			nil,
		),
		TorrentsActive: prometheus.NewDesc(
			namespace+collectorNamespace+"torrents_active",
			"The number of active torrents",
			nil,
			nil,
		),
		TorrentsPaused: prometheus.NewDesc(
			namespace+collectorNamespace+"torrents_paused",
			"The number of paused torrents",
			nil,
			nil,
		),

		Downloaded: prometheus.NewDesc(
			namespace+collectorNamespace+"downloaded_bytes",
			"The number of downloaded bytes",
			[]string{"type"},
			nil,
		),
		Uploaded: prometheus.NewDesc(
			namespace+collectorNamespace+"uploaded_bytes",
			"The number of uploaded bytes",
			[]string{"type"},
			nil,
		),
		FilesAdded: prometheus.NewDesc(
			namespace+collectorNamespace+"files_added",
			"The number of files added",
			[]string{"type"},
			nil,
		),
		ActiveTime: prometheus.NewDesc(
			namespace+collectorNamespace+"active",
			"The time transmission is active since",
			[]string{"type"},
			nil,
		),
		SessionCount: prometheus.NewDesc(
			namespace+collectorNamespace+"sessions",
			"Count of the times transmission started",
			[]string{"type"},
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
	stats, err := sc.client.GetSessionStats()
	if err != nil {
		log.Printf("failed to get session stats: %v", err)
	}

	dur := time.Duration(stats.CurrentStats.SecondsActive) * time.Second
	timestamp := time.Now().Add(-1 * dur).Unix()

	ch <- prometheus.MustNewConstMetric(
		sc.DownloadSpeed,
		prometheus.GaugeValue,
		float64(stats.DownloadSpeed),
	)
	ch <- prometheus.MustNewConstMetric(
		sc.UploadSpeed,
		prometheus.GaugeValue,
		float64(stats.UploadSpeed),
	)
	ch <- prometheus.MustNewConstMetric(
		sc.TorrentsTotal,
		prometheus.GaugeValue,
		float64(stats.TorrentCount),
	)
	ch <- prometheus.MustNewConstMetric(
		sc.TorrentsActive,
		prometheus.GaugeValue,
		float64(stats.ActiveTorrentCount),
	)
	ch <- prometheus.MustNewConstMetric(
		sc.TorrentsPaused,
		prometheus.GaugeValue,
		float64(stats.PausedTorrentCount),
	)

	types := []string{"current", "cumulative"}
	for _, t := range types {
		ch <- prometheus.MustNewConstMetric(
			sc.Downloaded,
			prometheus.GaugeValue,
			float64(stats.CurrentStats.DownloadedBytes),
			t,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.Uploaded,
			prometheus.GaugeValue,
			float64(stats.CurrentStats.UploadedBytes),
			t,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.FilesAdded,
			prometheus.GaugeValue,
			float64(stats.CurrentStats.FilesAdded),
			t,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.ActiveTime,
			prometheus.GaugeValue,
			float64(timestamp),
			t,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.SessionCount,
			prometheus.GaugeValue,
			float64(stats.CurrentStats.SessionCount),
			t,
		)
	}
}
