package apihop

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListModelsAndAuthHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Fatalf("unexpected auth header: %q", r.Header.Get("Authorization"))
		}
		if r.URL.Path != "/v1/models" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"object": "list", "data": []any{}})
	}))
	defer srv.Close()

	c := NewClient("test-key")
	c.BaseURL = srv.URL
	res, err := c.ListModels(context.Background())
	if err != nil {
		t.Fatalf("ListModels failed: %v", err)
	}
	if res["object"] != "list" {
		t.Fatalf("unexpected object: %v", res["object"])
	}
}

func TestTwitterUserByUsername(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/twitter/users" || r.URL.Query().Get("username") != "openai" {
			t.Fatalf("unexpected url: %s", r.URL.String())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"success": true})
	}))
	defer srv.Close()

	c := NewClient("test-key")
	c.BaseURL = srv.URL
	res, err := c.TwitterUserByUsername(context.Background(), "@openai")
	if err != nil {
		t.Fatalf("TwitterUserByUsername failed: %v", err)
	}
	if res["success"] != true {
		t.Fatalf("unexpected response: %v", res)
	}
}
