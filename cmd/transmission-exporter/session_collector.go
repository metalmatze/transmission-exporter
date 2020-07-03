package main

import (
	"log"

	"github.com/metalmatze/transmission-exporter"
	"github.com/prometheus/client_golang/prometheus"
)

// SessionCollector exposes session metrics
type SessionCollector struct {
	clients []*transmission.Client

	AltSpeedDown     *prometheus.Desc
	AltSpeedUp       *prometheus.Desc
	CacheSize        *prometheus.Desc
	FreeSpace        *prometheus.Desc
	QueueDown        *prometheus.Desc
	QueueUp          *prometheus.Desc
	PeerLimitGlobal  *prometheus.Desc
	PeerLimitTorrent *prometheus.Desc
	SeedRatioLimit   *prometheus.Desc
	SpeedLimitDown   *prometheus.Desc
	SpeedLimitUp     *prometheus.Desc
	Version          *prometheus.Desc
}

// NewSessionCollector takes a transmission.Client and returns a SessionCollector
func NewSessionCollector(clients []*transmission.Client) *SessionCollector {
	return &SessionCollector{
		clients: clients,

		AltSpeedDown: prometheus.NewDesc(
			namespace+"alt_speed_down",
			"Alternative max global download speed",
			[]string{"enabled", "client_name"},
			nil,
		),
		AltSpeedUp: prometheus.NewDesc(
			namespace+"alt_speed_up",
			"Alternative max global upload speed",
			[]string{"enabled", "client_name"},
			nil,
		),
		CacheSize: prometheus.NewDesc(
			namespace+"cache_size_bytes",
			"Maximum size of the disk cache",
			[]string{"client_name"},
			nil,
		),
		FreeSpace: prometheus.NewDesc(
			namespace+"free_space",
			"Free space left on device to download to",
			[]string{"download_dir", "incomplete_dir", "client_name"},
			nil,
		),
		QueueDown: prometheus.NewDesc(
			namespace+"queue_down",
			"Max number of torrents to download at once",
			[]string{"enabled", "client_name"},
			nil,
		),
		QueueUp: prometheus.NewDesc(
			namespace+"queue_up",
			"Max number of torrents to upload at once",
			[]string{"enabled", "client_name"},
			nil,
		),
		PeerLimitGlobal: prometheus.NewDesc(
			namespace+"global_peer_limit",
			"Maximum global number of peers",
			[]string{"client_name"},
			nil,
		),
		PeerLimitTorrent: prometheus.NewDesc(
			namespace+"torrent_peer_limit",
			"Maximum number of peers for a single torrent",
			[]string{"client_name"},
			nil,
		),
		SeedRatioLimit: prometheus.NewDesc(
			namespace+"seed_ratio_limit",
			"The default seed ratio for torrents to use",
			[]string{"enabled", "client_name"},
			nil,
		),
		SpeedLimitDown: prometheus.NewDesc(
			namespace+"speed_limit_down_bytes",
			"Max global download speed",
			[]string{"enabled", "client_name"},
			nil,
		),
		SpeedLimitUp: prometheus.NewDesc(
			namespace+"speed_limit_up_bytes",
			"Max global upload speed",
			[]string{"enabled", "client_name"},
			nil,
		),
		Version: prometheus.NewDesc(
			namespace+"version",
			"Transmission version as label",
			[]string{"version", "client_name"},
			nil,
		),
	}
}

// Describe implements the prometheus.Collector interface
func (sc *SessionCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- sc.AltSpeedDown
	ch <- sc.AltSpeedUp
	ch <- sc.CacheSize
	ch <- sc.FreeSpace
	ch <- sc.QueueDown
	ch <- sc.QueueUp
	ch <- sc.PeerLimitGlobal
	ch <- sc.PeerLimitTorrent
	ch <- sc.SeedRatioLimit
	ch <- sc.SpeedLimitDown
	ch <- sc.SpeedLimitUp
	ch <- sc.Version
}

// Collect implements the prometheus.Collector interface
func (sc *SessionCollector) Collect(ch chan<- prometheus.Metric) {
	for _, client := range sc.clients {

		session, err := client.GetSession()
		if err != nil {
			log.Printf("failed to get session: %v", err)
		}

		ch <- prometheus.MustNewConstMetric(
			sc.AltSpeedDown,
			prometheus.GaugeValue,
			float64(session.AltSpeedDown),
			boolToString(session.AltSpeedEnabled),
			client.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.AltSpeedUp,
			prometheus.GaugeValue,
			float64(session.AltSpeedUp),
			boolToString(session.AltSpeedEnabled),
			client.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.CacheSize,
			prometheus.GaugeValue,
			float64(session.CacheSizeMB*1024*1024),
			client.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.FreeSpace,
			prometheus.GaugeValue,
			float64(session.DownloadDirFreeSpace),
			session.DownloadDir, session.IncompleteDir, client.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.QueueDown,
			prometheus.GaugeValue,
			float64(session.DownloadQueueSize),
			boolToString(session.DownloadQueueEnabled),
			client.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.QueueUp,
			prometheus.GaugeValue,
			float64(session.SeedQueueSize),
			boolToString(session.SeedQueueEnabled),
			client.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.PeerLimitGlobal,
			prometheus.GaugeValue,
			float64(session.PeerLimitGlobal),
			client.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.PeerLimitTorrent,
			prometheus.GaugeValue,
			float64(session.PeerLimitPerTorrent),
			client.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.SeedRatioLimit,
			prometheus.GaugeValue,
			float64(session.SeedRatioLimit),
			boolToString(session.SeedRatioLimited),
			client.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.SpeedLimitDown,
			prometheus.GaugeValue,
			float64(session.SpeedLimitDown),
			boolToString(session.SpeedLimitDownEnabled),
			client.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.SpeedLimitUp,
			prometheus.GaugeValue,
			float64(session.SpeedLimitUp),
			boolToString(session.SpeedLimitUpEnabled),
			client.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			sc.Version,
			prometheus.GaugeValue,
			float64(1),
			session.Version,
			client.Name,
		)
	}
}
