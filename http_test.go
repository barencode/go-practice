package main

import (
	"encoding/json"
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

func TestStatusHandler(t *testing.T) {
	mine := Mine{Coal: 4, Reserves: 7, Closed: true}
	api := API{Mine: &mine}

	request := httptest.NewRequest(http.MethodGet, "/status", nil)
	response := httptest.NewRecorder()
	api.statusHandler(response, request)
	contentType := response.Header().Get("Content-Type")

	var got MineStatus
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&got); err != nil {
		t.Fatalf("decode response: %+v", err)
	}

	var want MineStatus
	want = MineStatus{Coal: 4, Reserves: 7, Closed: true}
	if got != want {
		t.Errorf("got status=%+v; want %+v", got, want)
	}

	if response.Code != http.StatusOK || contentType != "application/json" {
		t.Errorf("got code=%d and content type=%q; want code=200 and \"application/json\"", response.Code, contentType)
	}
}

func TestStatusHandlerMethodNotAllowed(t *testing.T) {
	mine := Mine{Coal: 4, Reserves: 7, Closed: true}
	api := API{Mine: &mine}

	request := httptest.NewRequest(http.MethodPost, "/status", nil)
	response := httptest.NewRecorder()
	api.statusHandler(response, request)

	if response.Code != http.StatusMethodNotAllowed || response.Header().Get("Allow") != http.MethodGet {
		t.Errorf("got response code=%d and response header=%q; want 405 and Allow=\"GET\"", response.Code, response.Header().Get("Allow"))
	}
}
