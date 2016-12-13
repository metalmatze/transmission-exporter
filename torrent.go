package transmission

type (
	// Torrent represents a transmission torrent
	Torrent struct {
		ID            int     `json:"id"`
		Name          string  `json:"name"`
		Status        int     `json:"status"`
		Added         int     `json:"addedDate"`
		LeftUntilDone int     `json:"leftUntilDone"`
		Eta           int     `json:"eta"`
		UploadRatio   float64 `json:"uploadRatio"`
		RateDownload  int     `json:"rateDownload"`
		RateUpload    int     `json:"rateUpload"`
		DownloadDir   string  `json:"downloadDir"`
		IsFinished    bool    `json:"isFinished"`
		PercentDone   float64 `json:"percentDone"`
		SeedRatioMode int     `json:"seedRatioMode"`
		HashString    string  `json:"hashString"`
		Error         int     `json:"error"`
		ErrorString   string  `json:"errorString"`
	}

	// ByID implements the sort Interface to sort by ID
	ByID []Torrent
	// ByName implements the sort Interface to sort by Name
	ByName []Torrent
	// ByDate implements the sort Interface to sort by Date
	ByDate []Torrent
	// ByRatio implements the sort Interface to sort by Ratio
	ByRatio []Torrent
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
