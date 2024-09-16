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

	Status             *prometheus.Desc
	Added              *prometheus.Desc
	Files              *prometheus.Desc
	Finished           *prometheus.Desc
	Done               *prometheus.Desc
	Ratio              *prometheus.Desc
	Download           *prometheus.Desc
	Upload             *prometheus.Desc
	PeersGettingFromUs *prometheus.Desc

	// TrackerStats
	Downloads *prometheus.Desc
	Leechers  *prometheus.Desc
	Seeders   *prometheus.Desc
}

// NewTorrentCollector creates a new torrent collector with the transmission.Client
func NewTorrentCollector(client *transmission.Client) *TorrentCollector {
	const collectorNamespace = "torrent_"

	return &TorrentCollector{
		client: client,

		Status: prometheus.NewDesc(
			namespace+collectorNamespace+"status",
			"Status of a torrent",
			[]string{"id", "name"},
			nil,
		),
		Added: prometheus.NewDesc(
			namespace+collectorNamespace+"added",
			"The unixtime time a torrent was added",
			[]string{"id", "name"},
			nil,
		),
		Files: prometheus.NewDesc(
			namespace+collectorNamespace+"files_total",
			"The total number of files in a torrent",
			[]string{"id", "name"},
			nil,
		),
		Finished: prometheus.NewDesc(
			namespace+collectorNamespace+"finished",
			"Indicates if a torrent is finished (1) or not (0)",
			[]string{"id", "name"},
			nil,
		),
		Done: prometheus.NewDesc(
			namespace+collectorNamespace+"done",
			"The percent of a torrent being done",
			[]string{"id", "name"},
			nil,
		),
		Ratio: prometheus.NewDesc(
			namespace+collectorNamespace+"ratio",
			"The upload ratio of a torrent",
			[]string{"id", "name"},
			nil,
		),
		Download: prometheus.NewDesc(
			namespace+collectorNamespace+"download_bytes",
			"The current download rate of a torrent in bytes",
			[]string{"id", "name"},
			nil,
		),
		Upload: prometheus.NewDesc(
			namespace+collectorNamespace+"upload_bytes",
			"The current upload rate of a torrent in bytes",
			[]string{"id", "name"},
			nil,
		),
		PeersGettingFromUs: prometheus.NewDesc(
			namespace+collectorNamespace+"peers_getting_from_us",
			"The current number of peers downloading from us",
			[]string{"id", "name"},
			nil,
		),

		// TrackerStats
		Downloads: prometheus.NewDesc(
			namespace+collectorNamespace+"downloads_total",
			"How often this torrent was downloaded",
			[]string{"id", "name", "tracker"},
			nil,
		),
		Leechers: prometheus.NewDesc(
			namespace+collectorNamespace+"leechers",
			"The number of peers downloading this torrent",
			[]string{"id", "name", "tracker"},
			nil,
		),
		Seeders: prometheus.NewDesc(
			namespace+collectorNamespace+"seeders",
			"The number of peers uploading this torrent",
			[]string{"id", "name", "tracker"},
			nil,
		),
	}
}

// Describe implements the prometheus.Collector interface
func (tc *TorrentCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- tc.Status
	ch <- tc.Added
	ch <- tc.Files
	ch <- tc.Finished
	ch <- tc.Done
	ch <- tc.Ratio
	ch <- tc.Download
	ch <- tc.Upload
	ch <- tc.Downloads
	ch <- tc.Leechers
	ch <- tc.Seeders
	ch <- tc.PeersGettingFromUs
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
		ch <- prometheus.MustNewConstMetric(
			tc.PeersGettingFromUs,
			prometheus.GaugeValue,
			float64(t.PeersGettingFromUs),
			id, t.Name,
		)

		tstats := make(map[string]transmission.TrackerStat)

		for _, tracker := range t.TrackerStats {
			if tr, exists := tstats[tracker.Host]; exists {
				tr.DownloadCount += tracker.DownloadCount
			} else {
				tstats[tracker.Host] = tracker
			}
		}

		for _, tracker := range tstats {
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
