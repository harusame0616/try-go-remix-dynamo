package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	srv := httptest.NewServer(NewMux())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	t.Run("ステータスコード200を返す", func(t *testing.T) {
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("JSON形式でstatus:okを返す", func(t *testing.T) {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}
		body := string(b)

		expected := `{"status":"ok"}`
		if body != expected {
			t.Errorf("expected body %q, got %q", expected, body)
		}
	})
}
