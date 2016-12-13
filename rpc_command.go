package transmission

type (
	// RPCCommand is the root command to interact with Transmission via RPC
	RPCCommand struct {
		Method    string       `json:"method,omitempty"`
		Arguments RPCArguments `json:"arguments,omitempty"`
		Result    string       `json:"result,omitempty"`
	}
	// RPCArguments specifies the RPCCommand in more detail
	RPCArguments struct {
		Fields       []string        `json:"fields,omitempty"`
		Torrents     []Torrent       `json:"torrents,omitempty"`
		Ids          []int           `json:"ids,omitempty"`
		DeleteData   bool            `json:"delete-local-data,omitempty"`
		DownloadDir  string          `json:"download-dir,omitempty"`
		MetaInfo     string          `json:"metainfo,omitempty"`
		Filename     string          `json:"filename,omitempty"`
		TorrentAdded RPCTorrentAdded `json:"torrent-added"`
	}
	// RPCTorrentAdded specifies the torrent to get added data from
	RPCTorrentAdded struct {
		HashString string `json:"hashString"`
		ID         int    `json:"id"`
		Name       string `json:"name"`
	}
)
