package telegram

import (
	"adviser-bot/lib/utils"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdates  = "getUpdates"
	sendMessage = "sendMessage"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (cl *Client) Updates(offset int, limit int) ([]Update, error) {
	const errorMsg = "can`t get updates"

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := cl.doRequest(getUpdates, q)
	if err != nil {
		return nil, utils.WrapError(errorMsg, err)
	}
	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, utils.WrapError(errorMsg, err)
	}

	return res.Result, nil
}
func (cl *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := cl.doRequest(sendMessage, q)
	if err != nil {
		return utils.WrapError("can`t send message", err)
	}
	return nil
}

func (cl *Client) doRequest(method string, query url.Values) ([]byte, error) {
	const errorMsg = "can`t do request"
	u := url.URL{
		Scheme: "https",
		Host:   cl.host,
		Path:   path.Join(cl.basePath, method),
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, utils.WrapError(errorMsg, err)
	}
	req.URL.RawQuery = query.Encode()
	resp, err := cl.client.Do(req)
	if err != nil {
		return nil, utils.WrapError(errorMsg, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.WrapError(errorMsg, err)
	}

	return body, nil
}
