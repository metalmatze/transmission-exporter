package transmission

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const endpoint = "/transmission/rpc/"

type (
	// User to authenticate with Transmission
	User struct {
		Username string
		Password string
	}
	// Client connects to transmission via HTTP
	Client struct {
		URL   string
		token string

		User   *User
		client http.Client
	}
)

// New create new transmission torrent
func New(url string, user *User) *Client {
	return &Client{
		URL:  url + endpoint,
		User: user,
	}
}

func (c *Client) post(body []byte) ([]byte, error) {
	authRequest, err := c.authRequest("POST", body)
	if err != nil {
		return make([]byte, 0), err
	}

	res, err := c.client.Do(authRequest)
	if err != nil {
		return make([]byte, 0), err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return make([]byte, 0), errors.New("authorization failed, check your username and password and make sure the ip is whitelisted")
	}

	if res.StatusCode == http.StatusConflict {
		c.getToken()
		authRequest, err := c.authRequest("POST", body)
		if err != nil {
			return make([]byte, 0), err
		}
		res, err = c.client.Do(authRequest)
		if err != nil {
			return make([]byte, 0), err
		}
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return make([]byte, 0), err
	}
	return resBody, nil
}

func (c *Client) getToken() error {
	req, err := http.NewRequest("POST", c.URL, strings.NewReader(""))
	if err != nil {
		return err
	}

	if c.User != nil {
		req.SetBasicAuth(c.User.Username, c.User.Password)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	c.token = res.Header.Get("X-Transmission-Session-Id")
	return nil
}

func (c *Client) authRequest(method string, body []byte) (*http.Request, error) {
	if c.token == "" {
		err := c.getToken()
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, c.URL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Transmission-Session-Id", c.token)

	if c.User != nil {
		req.SetBasicAuth(c.User.Username, c.User.Password)
	}

	return req, nil
}

// GetTorrents get a list of torrents
func (c *Client) GetTorrents() ([]Torrent, error) {
	cmd := TorrentCommand{
		Method: "torrent-get",
		Arguments: TorrentArguments{
			Fields: []string{
				"id",
				"name",
				"hashString",
				"status",
				"addedDate",
				"leftUntilDone",
				"eta",
				"uploadRatio",
				"rateDownload",
				"rateUpload",
				"downloadDir",
				"isFinished",
				"percentDone",
				"seedRatioMode",
				"error",
				"errorString",
				"files",
				"fileStats",
				"peers",
				"trackers",
				"trackerStats",
				"peersGettingFromUs",
				"totalSize",
			},
		},
	}

	req, err := json.Marshal(&cmd)
	if err != nil {
		return nil, err
	}

	resp, err := c.post(req)
	if err != nil {
		return nil, err
	}

	var out TorrentCommand
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}

	return out.Arguments.Torrents, nil
}

// GetSession gets the current session from transmission
func (c *Client) GetSession() (*Session, error) {
	req, err := json.Marshal(SessionCommand{Method: "session-get"})
	if err != nil {
		return nil, err
	}

	resp, err := c.post(req)
	if err != nil {
		return nil, err
	}

	var cmd SessionCommand
	if err := json.Unmarshal(resp, &cmd); err != nil {
		return nil, err
	}

	return &cmd.Session, nil
}

// GetSessionStats gets stats on the current & cumulative session
func (c *Client) GetSessionStats() (*SessionStats, error) {
	req, err := json.Marshal(SessionCommand{Method: "session-stats"})
	if err != nil {
		return nil, err
	}

	resp, err := c.post(req)
	if err != nil {
		return nil, err
	}

	var cmd SessionStatsCmd
	if err := json.Unmarshal(resp, &cmd); err != nil {
		return nil, err
	}

	return &cmd.SessionStats, nil
}
