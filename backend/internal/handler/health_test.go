package handler

import (
	"net/http"
	"testing"
)

func TestHealthCheck_Success(t *testing.T) {
	r := newRouter(0)
	r.GET("/health", HealthCheck)
	w := doRequest(r, "GET", "/health", nil)

	assertStatus(t, w, http.StatusOK)
	body := parseJSON(t, w)
	if status, ok := body["status"].(string); !ok || status != "ok" {
		t.Errorf("expected status 'ok', got %v", body["status"])
	}
}
