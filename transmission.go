package transmission

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

const endpoint = "/transmission/rpc"

type (
	User struct {
		Username string
		Password string
	}
	Client struct {
		URL   string
		token string

		User   *User
		client http.Client
	}
)

//New create new transmission torrent
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

	if res.StatusCode == 409 {
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

func (c *Client) ExecuteCommand(cmd *RPCCommand) (*RPCCommand, error) {
	var out RPCCommand

	body, err := json.Marshal(&cmd)
	if err != nil {
		return nil, err
	}

	output, err := c.post(body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(output, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

//GetTorrents get a list of torrents
func (ac *Client) GetTorrents() ([]Torrent, error) {
	cmd := &RPCCommand{
		Method: "torrent-get",
		Arguments: RPCArguments{
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
			},
		},
	}

	out, err := ac.ExecuteCommand(cmd)
	if err != nil {
		return nil, err
	}

	return out.Arguments.Torrents, nil
}
