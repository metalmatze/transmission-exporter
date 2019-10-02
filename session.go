package transmission

type (
	// SessionCommand is the root command to interact with Transmission via RPC
	SessionCommand struct {
		Method  string  `json:"method,omitempty"`
		Session Session `json:"arguments,omitempty"`
		Result  string  `json:"result,omitempty"`
	}

	// Session information about the current transmission session
	Session struct {
		AltSpeedDown    int  `json:"alt-speed-down"`
		AltSpeedEnabled bool `json:"alt-speed-enabled"`
		//Alt_speed_time_begin   int    `json:"alt-speed-time-begin"`
		//Alt_speed_time_day     int    `json:"alt-speed-time-day"`
		//Alt_speed_time_enabled bool   `json:"alt-speed-time-enabled"`
		//Alt_speed_time_end     int    `json:"alt-speed-time-end"`
		AltSpeedUp int `json:"alt-speed-up"`
		//Blocklist_enabled       bool   `json:"blocklist-enabled"`
		//Blocklist_size          int    `json:"blocklist-size"`
		//Blocklist_url           string `json:"blocklist-url"`
		CacheSizeMB int `json:"cache-size-mb"`
		//Config_dir                   string `json:"config-dir"`
		//Dht_enabled                  bool   `json:"dht-enabled"`
		DownloadDir          string `json:"download-dir"`
		DownloadDirFreeSpace int64  `json:"download-dir-free-space"`
		DownloadQueueEnabled bool   `json:"download-queue-enabled"`
		DownloadQueueSize    int    `json:"download-queue-size"`
		//Encryption                 string `json:"encryption"`
		//Idle_seeding_limit           int    `json:"idle-seeding-limit"`
		//Idle_seeding_limit_enabled   bool   `json:"idle-seeding-limit-enabled"`
		IncompleteDir string `json:"incomplete-dir"`
		//Incomplete_dir_enabled       bool   `json:"incomplete-dir-enabled"`
		//Lpd_enabled                  bool   `json:"lpd-enabled"`
		PeerLimitGlobal     int `json:"peer-limit-global"`
		PeerLimitPerTorrent int `json:"peer-limit-per-torrent"`
		//Peer_port                 int    `json:"peer-port"`
		//Peer_port_random_on_start bool   `json:"peer-port-random-on-start"`
		//Pex_enabled               bool   `json:"pex-enabled"`
		//Port_forwarding_enabled   bool   `json:"port-forwarding-enabled"`
		//Queue_stalled_enabled     bool   `json:"queue-stalled-enabled"`
		//Queue_stalled_minutes     int    `json:"queue-stalled-minutes"`
		//Rename_partial_files      bool   `json:"rename-partial-files"`
		//RPC_version               int    `json:"rpc-version"`
		//RPC_version_minimum          int    `json:"rpc-version-minimum"`
		//Script_torrent_done_enabled  bool   `json:"script-torrent-done-enabled"`
		//Script_torrent_done_filename string `json:"script-torrent-done-filename"`
		SeedQueueEnabled      bool    `json:"seed-queue-enabled"`
		SeedQueueSize         int     `json:"seed-queue-size"`
		SeedRatioLimit        float64 `json:"seedRatioLimit"`
		SeedRatioLimited      bool    `json:"seedRatioLimited"`
		SpeedLimitDown        int     `json:"speed-limit-down"`
		SpeedLimitDownEnabled bool    `json:"speed-limit-down-enabled"`
		SpeedLimitUp          int     `json:"speed-limit-up"`
		SpeedLimitUpEnabled   bool    `json:"speed-limit-up-enabled"`
		//Start_added_torrents         bool   `json:"start-added-torrents"`
		//Trash_original_torrent_files bool   `json:"trash-original-torrent-files"`
		//Utp_enabled                  bool   `json:"utp-enabled"`
		Version string `json:"version"`
	}
)
