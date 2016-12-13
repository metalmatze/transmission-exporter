package transmission

type (
	RPCCommand struct {
		Method    string       `json:"method,omitempty"`
		Arguments RPCArguments `json:"arguments,omitempty"`
		Result    string       `json:"result,omitempty"`
	}
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
	RPCTorrentAdded struct {
		HashString string `json:"hashString"`
		ID         int    `json:"id"`
		Name       string `json:"name"`
	}
)
