package transmission

type (
	// SessionStatsCmd is the root command to interact with Transmission via RPC
	SessionStatsCmd struct {
		SessionStats `json:"arguments"`
		Result       string `json:"result"`
	}

	// SessionStats contains information about the current & cumulative session
	SessionStats struct {
		DownloadSpeed      int64             `json:"downloadSpeed"`
		UploadSpeed        int64             `json:"uploadSpeed"`
		ActiveTorrentCount int               `json:"activeTorrentCount"`
		PausedTorrentCount int               `json:"pausedTorrentCount"`
		TorrentCount       int               `json:"torrentCount"`
		CumulativeStats    SessionStateStats `json:"cumulative-stats"`
		CurrentStats       SessionStateStats `json:"current-stats"`
	}
	// SessionStateStats contains current or cumulative session stats
	SessionStateStats struct {
		DownloadedBytes int64 `json:"downloadedBytes"`
		UploadedBytes   int64 `json:"uploadedBytes"`
		FilesAdded      int64 `json:"filesAdded"`
		SecondsActive   int64 `json:"secondsActive"`
		SessionCount    int64 `json:"sessionCount"`
	}
)
