package whatsappworkerclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientSendAndStatus(t *testing.T) {
	var gotSend SendRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/send":
			if r.Method != http.MethodPost {
				t.Fatalf("method: %s", r.Method)
			}
			body, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(body, &gotSend); err != nil {
				t.Fatal(err)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok":true}`))
		case "/v1/status":
			_, _ = w.Write([]byte(`{"connected":true}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL})
	if !c.Enabled() {
		t.Fatal("expected enabled client")
	}

	st, err := c.Status(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !st.Connected {
		t.Fatal("expected connected")
	}

	err = c.Send(context.Background(), SendRequest{
		Message: "hello",
		Type:    "problem",
		VideoID: "vid-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotSend.Message != "hello" || gotSend.Type != "problem" {
		t.Fatalf("unexpected payload: %+v", gotSend)
	}
}
