package transmission

type (
	// SessionStatsCmd is the root command to interact with Transmission via RPC
	SessionStatsCmd struct {
		SessionStats `json:"arguments"`
		Result       string `json:"result"`
	}

	// SessionStats contains information about the current & cumulative session
	SessionStats struct {
		DownloadSpeed      int               `json:"downloadSpeed"`
		UploadSpeed        int               `json:"uploadSpeed"`
		ActiveTorrentCount int               `json:"activeTorrentCount"`
		PausedTorrentCount int               `json:"pausedTorrentCount"`
		TorrentCount       int               `json:"torrentCount"`
		CumulativeStats    SessionStateStats `json:"cumulative-stats"`
		CurrentStats       SessionStateStats `json:"current-stats"`
	}
	// SessionStateStats contains current or cumulative session stats
	SessionStateStats struct {
		DownloadedBytes int `json:"downloadedBytes"`
		UploadedBytes   int `json:"uploadedBytes"`
		FilesAdded      int `json:"filesAdded"`
		SecondsActive   int `json:"secondsActive"`
		SessionCount    int `json:"sessionCount"`
	}
)
