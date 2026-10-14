package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	healthHandler(response, request)

	if response.Code != http.StatusOK || response.Body.String() != "Шахта на связи\n" {
		t.Errorf("got code=%d, body=%q; want code=200, body=%q", response.Code, response.Body.String(), "Шахта на связи\n")
	}
}
