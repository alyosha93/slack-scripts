package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/slack-go/slack"
)

func TestPostMsg(t *testing.T) {
	testCases := []struct {
		description   string
		msg           Msg
		respPostMsg   []byte
		wantTS        string
		wantChannelID string
		wantErr       string
	}{
		{
			description:   "successfully posted message",
			msg:           Msg{Body: "Hey!"},
			respPostMsg:   []byte(mockPostMsgResp),
			wantTS:        "1503435956.000247",
			wantChannelID: "C1H9RESGL",
		},
		{
			description: "failure to post message",
			msg:         Msg{Body: "Hey!"},
			respPostMsg: []byte(mockPostMsgErrResp),
			wantErr:     "too_many_attachments",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respPostMsg)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL)))

			ts, err := PostMsg(client, Msg{}, "C1H9RESGL")

			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantErr != "" {
				if err == nil {
					t.Fatal("expected error but did not receive one")
				}
				if err.Error() != tc.wantErr {
					t.Fatalf("expected to receive error: %s, got: %s", tc.wantErr, err)
				}
			}

			if ts != tc.wantTS {
				t.Fatalf("expected timestamp: %s, got: %s", tc.wantTS, ts)
			}
		})
	}
}

func TestPostThreadMsg(t *testing.T) {
	testCases := []struct {
		description string
		msg         Msg
		respPostMsg []byte
		wantErr     string
	}{
		{
			description: "successfully posted message",
			msg:         Msg{Body: "Hey!"},
			respPostMsg: []byte(mockPostMsgResp),
		},
		{
			description: "failure to post message",
			msg:         Msg{Body: "Hey!"},
			respPostMsg: []byte(mockPostMsgErrResp),
			wantErr:     "too_many_attachments",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/chat.postMessage", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respPostMsg)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL)))

			err := PostThreadMsg(client, Msg{}, "C1H9RESGL", "1503435956.000247")

			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantErr != "" {
				if err == nil {
					t.Fatal("expected error but did not receive one")
				}
				if err.Error() != tc.wantErr {
					t.Fatalf("expected to receive error: %s, got: %s", tc.wantErr, err)
				}
			}
		})
	}
}

func TestPostEphemeralMsg(t *testing.T) {
	testCases := []struct {
		description   string
		msg           Msg
		respPostMsg   []byte
		wantTS        string
		wantChannelID string
		wantErr       string
	}{
		{
			description:   "successfully posted ephemeral message",
			msg:           Msg{Body: "Hey!"},
			respPostMsg:   []byte(mockPostMsgResp),
			wantTS:        "1503435956.000247",
			wantChannelID: "C1H9RESGL",
		},
		{
			description: "failure to post ephemeral message",
			msg:         Msg{Body: "Hey!"},
			respPostMsg: []byte(mockPostMsgErrResp),
			wantErr:     "too_many_attachments",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/chat.postEphemeral", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respPostMsg)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL)))

			err := PostEphemeralMsg(client, Msg{}, "C1H9RESGL", "U12345")

			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantErr != "" {
				if err == nil {
					t.Fatal("expected error but did not receive one")
				}
				if err.Error() != tc.wantErr {
					t.Fatalf("expected to receive error: %s, got: %s", tc.wantErr, err)
				}
			}
		})
	}
}

func TestUpdateMsg(t *testing.T) {
	testCases := []struct {
		description   string
		msg           Msg
		respUpdateMsg []byte
		wantErr       string
	}{
		{
			description:   "successfully posted message",
			msg:           Msg{Body: "Hey!"},
			respUpdateMsg: []byte(mockUpdateMsgResp),
		},
		{
			description:   "failure to post message",
			msg:           Msg{Body: "Hey!"},
			respUpdateMsg: []byte(mockPostMsgErrResp),
			wantErr:       "too_many_attachments",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/chat.update", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(tc.respUpdateMsg)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL)))

			err := UpdateMsg(client, Msg{}, "C1H9RESGL", "1503435957.000237")

			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantErr != "" {
				if err == nil {
					t.Fatal("expected error but did not receive one")
				}
				if err.Error() != tc.wantErr {
					t.Fatalf("expected to receive error: %s, got: %s", tc.wantErr, err)
				}
			}
		})
	}
}

func TestDeleteMsg(t *testing.T) {
	testCases := []struct {
		description   string
		respDeleteMsg []byte
		wantErr       string
	}{
		{
			description:   "successfully posted message",
			respDeleteMsg: []byte(mockSuccessResp),
		},
		{
			description:   "failure to post message",
			respDeleteMsg: []byte(mockPostMsgErrResp),
			wantErr:       "too_many_attachments",
		},
	}

	responseURLPath := "/commands/XXXXXXXX/00000000/YYYYYYYY"
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(responseURLPath, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(tc.respDeleteMsg)
			})

			testServ := httptest.NewServer(mux)
			defer testServ.Close()

			client := slack.New("x012345", slack.OptionAPIURL(fmt.Sprintf("%v/", testServ.URL)))
			err := DeleteMsg(client, "C1H9RESGL", "1503435957.000237", fmt.Sprintf("%s%s", testServ.URL, responseURLPath))

			if tc.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantErr != "" {
				if err == nil {
					t.Fatal("expected error but did not receive one")
				}
				if err.Error() != tc.wantErr {
					t.Fatalf("expected to receive error: %s, got: %s", tc.wantErr, err)
				}
			}
		})
	}
}

func TestSendResp(t *testing.T) {
	var msg slack.Message
	handler := func(w http.ResponseWriter, r *http.Request) {
		err := SendResp(w, slack.Message{})
		if err != nil {
			t.Fatalf("unexpected error handing request: %s", err)
		}
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got: %v", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("content type application/json, got: %v", resp.Header.Get("Content-Type"))
	}

	err := json.Unmarshal(body, &msg)
	if err != nil {
		t.Fatalf("failed to unmarshal response with error: %s", err)
	}

	if msg.ReplaceOriginal {
		t.Fatal("replace original should be false, but is true")
	}

	if msg.DeleteOriginal {
		t.Fatal("delete original should be false, but is true")
	}
}

func TestReplaceOriginal(t *testing.T) {
	var msg slack.Message
	handler := func(w http.ResponseWriter, r *http.Request) {
		err := ReplaceOriginal(w, slack.Message{})
		if err != nil {
			t.Fatalf("unexpected error handing request: %s", err)
		}
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got: %v", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("content type application/json, got: %v", resp.Header.Get("Content-Type"))
	}

	err := json.Unmarshal(body, &msg)
	if err != nil {
		t.Fatalf("failed to unmarshal response with error: %s", err)
	}

	if !msg.ReplaceOriginal {
		t.Fatal("replace original should be true, but is false")
	}

	if msg.DeleteOriginal {
		t.Fatal("delete original should be false, but is true")
	}
}

func TestSendOKAndDeleteOriginal(t *testing.T) {
	var msg slack.Message
	handler := func(w http.ResponseWriter, r *http.Request) {
		err := SendOKAndDeleteOriginal(w)
		if err != nil {
			t.Fatalf("unexpected error handing request: %s", err)
		}
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got: %v", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("content type application/json, got: %v", resp.Header.Get("Content-Type"))
	}

	err := json.Unmarshal(body, &msg)
	if err != nil {
		t.Fatalf("failed to unmarshal response with error: %s", err)
	}

	if msg.ReplaceOriginal {
		t.Fatal("replace original should be false, but is true")
	}

	if !msg.DeleteOriginal {
		t.Fatal("delete original should be true, but is false")
	}
}

func TestSendEmptyOK(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		SendEmptyOK(w)
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got: %v", resp.StatusCode)
	}
}
