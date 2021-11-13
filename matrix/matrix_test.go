package matrix

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	expected := &Client{
		homeserver: "https://matrix.example.com",
		login:      "test",
		password:   "test",
		Room:       "!test:example.com",
		MsgType:    "m.text",
	}

	actual := New("https://matrix.example.com", "test", "test", "", "!test:example.com", "")

	if !reflect.DeepEqual(expected, actual) {
		t.Fail()
	}
}

func TestNew_Token(t *testing.T) {
	expected := &Client{
		homeserver: "https://matrix.example.com",
		token:      "test",
		nologin:    true,
		Room:       "!test:example.com",
		MsgType:    "m.notice",
	}

	actual := New("https://matrix.example.com", "", "", "test", "!test:example.com", "m.notice")

	if !reflect.DeepEqual(expected, actual) {
		t.Fail()
	}
}

func TestResolveRoom(t *testing.T) {
	ctx := context.TODO()
	response := `{"room_id":"!test:example.com","servers":["example.com"]}`
	client, server := startServer(t, "/_matrix/client/r0/directory/room/%23test:example.com", nil, []byte(response))
	defer server.Close()
	client.Room = "#test:example.com"

	roomID, err := client.ResolveRoom(ctx, "#test:example.com")
	if err != nil {
		t.Error(err)
	}
	if roomID != "!test:example.com" {
		t.Fail()
	}
}

func TestResolveRoom_NoAlias(t *testing.T) {
	ctx := context.TODO()
	client := &Client{
		homeserver: "https://matrix.example.com",
		Room:       "!test:example.com",
	}

	roomID, err := client.ResolveRoom(ctx, client.Room)
	if err != nil {
		t.Error(err)
	}
	if roomID != "!test:example.com" {
		t.Fail()
	}
}

func TestLogin(t *testing.T) {
	ctx := context.TODO()
	expectedRequest := `{"type":"m.login.password","identifier":{"type":"m.id.user","user":"test"},"password":"test"}`
	response := `{"access_token": "token","device_id": "deviceID","user_id": "@test:example.com"}`
	client, server := startServer(t, "/_matrix/client/r0/login", []byte(expectedRequest), []byte(response))
	defer server.Close()

	err := client.Login(ctx)
	if err != nil {
		t.Error(err)
	}
	if client.token != "token" {
		t.Error("incorrect token is set")
	}
}

func TestLogin_NoLogin(t *testing.T) {
	ctx := context.TODO()
	client := &Client{
		homeserver: "https://matrix.example.com",
		token:      "test",
		nologin:    true,
		Room:       "!test:example.com",
	}

	err := client.Login(ctx)
	if err != nil {
		t.Error(err)
	}
	if client.token != "test" {
		t.Error("incorrect token is set")
	}
}

func TestLogout(t *testing.T) {
	ctx := context.TODO()
	client, server := startServer(t, "/_matrix/client/r0/logout?access_token=test", nil, nil)
	client.token = "test"
	defer server.Close()

	client.logout(ctx)
}

func TestLogout_NoLogin(_ *testing.T) {
	// Not an actual test, because client.logout will just call "return" without any params in case of token auth
	ctx := context.TODO()
	client := &Client{
		homeserver: "https://matrix.example.com",
		token:      "test",
		nologin:    true,
		Room:       "!test:example.com",
	}

	client.logout(ctx)
}

func TestSendMessage(t *testing.T) {
	expectedPath := "/_matrix/client/r0/rooms/%21test:example.com/send/m.room.message?access_token=test"
	expectedRequestBody := `{"body":"hello!","formatted_body":"\u003cb\u003ehello!\u003c/b\u003e","format":"org.matrix.custom.html","msgtype":"m.text"}`
	responseBody := `{"access_token": "token","device_id": "deviceID","user_id": "@test:example.com"}`
	client := &Client{
		homeserver: "https://matrix.example.com",
		login:      "test",
		password:   "test",
		Room:       "!test:example.com",
		token:      "test",
		MsgType:    "m.text",
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualPath := r.URL.RequestURI()
		// cheating - we're trying to test only message sending here, so ignore other checks, because they are done in separate test cases
		if !strings.HasPrefix(actualPath, "/_matrix/client/r0/rooms") {
			w.Write(nil)
			return
		}
		if actualPath != expectedPath {
			t.Error("request url is not expected. actual:", actualPath)
		}
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error("sent request body cannot be read")
		}
		if !bytes.Equal([]byte(expectedRequestBody), requestBody) {
			t.Error("sent request body is not expected")
		}
		defer r.Body.Close()

		w.Write([]byte(responseBody))
	}))
	client.homeserver = server.URL
	defer server.Close()

	err := client.SendMessage("hello!", "<b>hello!</b>")
	if err != nil {
		t.Error(err)
	}
}

func startServer(t *testing.T, expectedPath string, expectedRequestBody []byte, responseBody []byte) (*Client, *httptest.Server) {
	t.Helper()
	client := &Client{
		homeserver: "https://matrix.example.com",
		login:      "test",
		password:   "test",
		Room:       "!test:example.com",
		MsgType:    "m.text",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualPath := r.URL.RequestURI()
		if actualPath != expectedPath {
			t.Error("request url is not expected. actual:", actualPath)
		}
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error("sent request body cannot be read")
		}
		if !bytes.Equal(expectedRequestBody, requestBody) {
			t.Log(string(expectedRequestBody), string(requestBody))
			t.Error("sent request body is not expected")
		}
		defer r.Body.Close()

		w.Write(responseBody)
	}))
	client.homeserver = server.URL

	return client, server
}
