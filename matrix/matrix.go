package matrix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// MaxPayloadSize of any matrix event, in bytes
const MaxPayloadSize = 65536

// InfrastructurePayloadSize is an approximate size in bytes of json payload template, may differ a bit by called command size
const InfrastructurePayloadSize = 1000

// SuggestedPayloadBuffer is a "just in case" spare size in bytes, suggested to leave unused
const SuggestedPayloadBuffer = 2000

// Client implementation
type Client struct {
	homeserver string
	password   string
	Room       string
	login      string
	token      string
	nologin    bool
	MsgType    string
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
	FormattedBody string `json:"formatted_body,omitempty"`
	Format        string `json:"format,omitempty"`
	MsgType       string `json:"msgtype"`
}

type roomResponse struct {
	RoomID string `json:"room_id"`
}

// New matrix client
func New(homeserver, login, password, token, room, msgtype string) *Client {
	var nologin bool
	if token != "" {
		nologin = true
	}
	if msgtype == "" {
		msgtype = "m.text"
	}

	return &Client{
		homeserver: homeserver,
		password:   password,
		nologin:    nologin,
		login:      login,
		token:      token,

		Room:    room,
		MsgType: msgtype,
	}
}

// Login as matrix user
func (c *Client) Login(ctx context.Context) error {
	if c.nologin {
		return nil
	}

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
	body, err := io.ReadAll(resp.Body)
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

// ResolveRoom returns room ID from alias
func (c *Client) ResolveRoom(ctx context.Context, room string) (string, error) {
	if !c.isRoomAlias() {
		return c.Room, nil
	}

	endpoint := fmt.Sprintf("%s/_matrix/client/r0/directory/room/%s", c.homeserver, url.PathEscape(room))
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return room, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return room, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return room, err
	}
	if resp.StatusCode != 200 {
		return room, fmt.Errorf("could not resolve room alias: %s", string(body))
	}

	var data roomResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return room, err
	}

	return data.RoomID, nil
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

// isRoomAlias checks if provided room is an alias or room id
func (c *Client) isRoomAlias() bool {
	return strings.HasPrefix(c.Room, "#")
}

func (c *Client) send(ctx context.Context, plaintext string, html string) error {
	var format string
	endpoint := fmt.Sprintf("%s/_matrix/client/r0/rooms/%s/send/m.room.message?access_token=%s", c.homeserver, url.PathEscape(c.Room), c.token)
	if html != "" {
		format = "org.matrix.custom.html"
	}
	request, err := json.Marshal(&messageRequest{
		Body:          plaintext,
		FormattedBody: html,
		Format:        format,
		MsgType:       c.MsgType,
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
	body, err := io.ReadAll(resp.Body)
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
	if c.nologin {
		return
	}

	endpoint := fmt.Sprintf("%s/_matrix/client/r0/logout?access_token=%s", c.homeserver, c.token)
	req, _ := http.NewRequestWithContext(ctx, "POST", endpoint, nil)
	http.DefaultClient.Do(req)
}
