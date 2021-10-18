package matrix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client implementation
type Client struct {
	homeserver string
	password   string
	roomID     string
	login      string
	token      string
}

type loginRequest struct {
	Type       string                 `json:"type"`
	Identifier loginRequestIdentifier `json:"identifier"`
	Password   string                 `json:"password"`
}

type loginRequestIdentifier struct {
	Type string `json:"type"`
	User string `json:"user"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

type messageRequest struct {
	Body          string `json:"body"`
	FormattedBody string `json:"formatted_body"`
	Format        string `json:"format"`
	MsgType       string `json:"msgtype"`
}

// New matrix client
func New(homeserver, login, password, roomID string) *Client {
	return &Client{
		homeserver: homeserver,
		password:   password,
		roomID:     roomID,
		login:      login,
	}
}

// Login as matrix user
func (c *Client) Login(ctx context.Context) error {
	endpoint := c.homeserver + "/_matrix/client/r0/login"
	request, err := json.Marshal(&loginRequest{
		Type: "m.login.password",
		Identifier: loginRequestIdentifier{
			Type: "m.id.user",
			User: c.login,
		},
		Password: c.password,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(request))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	response := &loginResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	c.token = response.AccessToken

	return nil
}

func (c *Client) send(ctx context.Context, plaintext string, html string) error {
	endpoint := fmt.Sprintf("%s/_matrix/client/r0/rooms/%s/send/m.room.message?access_token=%s", c.homeserver, url.PathEscape(c.roomID), c.token)
	request, err := json.Marshal(&messageRequest{
		Body:          plaintext,
		FormattedBody: html,
		Format:        "org.matrix.custom.html",
		MsgType:       "m.text",
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(request))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("message was not sent: %s", string(body))
	}
	return nil
}

// nolint // nobody cares about error here, worst case - the session will not be destroyed
func (c *Client) logout(ctx context.Context) {
	endpoint := fmt.Sprintf("%s/_matrix/client/r0/logout?access_token=%s", c.homeserver, c.token)
	req, _ := http.NewRequestWithContext(ctx, "POST", endpoint, nil)
	http.DefaultClient.Do(req)
}

// SendMessage to the matrix room.
func (c *Client) SendMessage(plaintext, html string) error {
	// that cycle needed in case login() goroutine in main package didn't finish yet
	for {
		if c.token != "" {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	ctx := context.Background()
	defer c.logout(ctx)

	err := c.send(ctx, plaintext, html)
	if err != nil {
		return err
	}

	return nil
}
