package transmission

type (
	// Torrent represents a transmission torrent
	Torrent struct {
		ID            int           `json:"id"`
		Name          string        `json:"name"`
		Status        int           `json:"status"`
		Added         int           `json:"addedDate"`
		LeftUntilDone int           `json:"leftUntilDone"`
		Eta           int           `json:"eta"`
		UploadRatio   float64       `json:"uploadRatio"`
		RateDownload  int           `json:"rateDownload"`
		RateUpload    int           `json:"rateUpload"`
		DownloadDir   string        `json:"downloadDir"`
		IsFinished    bool          `json:"isFinished"`
		PercentDone   float64       `json:"percentDone"`
		SeedRatioMode int           `json:"seedRatioMode"`
		HashString    string        `json:"hashString"`
		Error         int           `json:"error"`
		ErrorString   string        `json:"errorString"`
		Files         []File        `json:"files"`
		FilesStats    []FileStat    `json:"fileStats"`
		TrackerStats  []TrackerStat `json:"trackerStats"`
		Peers         []Peer        `json:"peers"`
	}

	// ByID implements the sort Interface to sort by ID
	ByID []Torrent
	// ByName implements the sort Interface to sort by Name
	ByName []Torrent
	// ByDate implements the sort Interface to sort by Date
	ByDate []Torrent
	// ByRatio implements the sort Interface to sort by Ratio
	ByRatio []Torrent

	// File is a file contained inside a torrent
	File struct {
		BytesCompleted int    `json:"bytesCompleted"`
		Length         int    `json:"length"`
		Name           string `json:"name"`
	}

	// FileStat describe a file's priority & if it's wanted
	FileStat struct {
		BytesCompleted int  `json:"bytesCompleted"`
		Priority       int  `json:"priority"`
		Wanted         bool `json:"wanted"`
	}

	// TrackerStat has stats about the torrent's tracker
	TrackerStat struct {
		Announce              string `json:"announce"`
		AnnounceState         int    `json:"announceState"`
		DownloadCount         int    `json:"downloadCount"`
		HasAnnounced          bool   `json:"hasAnnounced"`
		HasScraped            bool   `json:"hasScraped"`
		Host                  string `json:"host"`
		ID                    int    `json:"id"`
		IsBackup              bool   `json:"isBackup"`
		LastAnnouncePeerCount int    `json:"lastAnnouncePeerCount"`
		LastAnnounceResult    string `json:"lastAnnounceResult"`
		LastAnnounceStartTime int    `json:"lastAnnounceStartTime"`
		LastAnnounceSucceeded bool   `json:"lastAnnounceSucceeded"`
		LastAnnounceTime      int    `json:"lastAnnounceTime"`
		LastAnnounceTimedOut  bool   `json:"lastAnnounceTimedOut"`
		LastScrapeResult      string `json:"lastScrapeResult"`
		LastScrapeStartTime   int    `json:"lastScrapeStartTime"`
		LastScrapeSucceeded   bool   `json:"lastScrapeSucceeded"`
		LastScrapeTime        int    `json:"lastScrapeTime"`
		LastScrapeTimedOut    int    `json:"lastScrapeTimedOut"`
		LeecherCount          int    `json:"leecherCount"`
		NextAnnounceTime      int    `json:"nextAnnounceTime"`
		NextScrapeTime        int    `json:"nextScrapeTime"`
		Scrape                string `json:"scrape"`
		ScrapeState           int    `json:"scrapeState"`
		SeederCount           int    `json:"seederCount"`
		Tier                  int    `json:"tier"`
	}

	// Peer of a torrent
	Peer struct {
		Address            string `json:"address"`
		ClientIsChoked     bool   `json:"clientIsChoked"`
		ClientIsInterested bool   `json:"clientIsInterested"`
		ClientName         string `json:"clientName"`
		FlagStr            string `json:"flagStr"`
		IsDownloadingFrom  bool   `json:"isDownloadingFrom"`
		IsEncrypted        bool   `json:"isEncrypted"`
		IsIncoming         bool   `json:"isIncoming"`
		IsUTP              bool   `json:"isUTP"`
		IsUploadingTo      bool   `json:"isUploadingTo"`
		PeerIsChoked       bool   `json:"peerIsChoked"`
		PeerIsInterested   bool   `json:"peerIsInterested"`
		Port               int    `json:"port"`
		Progress           int    `json:"progress"`
		RateToClient       int    `json:"rateToClient"`
		RateToPeer         int    `json:"rateToPeer"`
	}
)

func (t ByID) Len() int           { return len(t) }
func (t ByID) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByID) Less(i, j int) bool { return t[i].ID < t[j].ID }

func (t ByName) Len() int           { return len(t) }
func (t ByName) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByName) Less(i, j int) bool { return t[i].Name < t[j].Name }

func (t ByDate) Len() int           { return len(t) }
func (t ByDate) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByDate) Less(i, j int) bool { return t[i].Added < t[j].Added }

func (t ByRatio) Len() int           { return len(t) }
func (t ByRatio) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByRatio) Less(i, j int) bool { return t[i].UploadRatio < t[j].UploadRatio }
