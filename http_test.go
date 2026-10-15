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

func TestMinerHandlerFound(t *testing.T) {
	mine := Mine{Coal: 4, Reserves: 7, Closed: true}
	mine.AddMiner(101, &Miner{Name: "Борис", Energy: 2, Production: 3})
	api := API{Mine: &mine}

	request := httptest.NewRequest(http.MethodGet, "/miner?id=101", nil)
	response := httptest.NewRecorder()
	api.minerHandler(response, request)
	contentType := response.Header().Get("Content-Type")

	var got Miner
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&got); err != nil {
		t.Fatalf("decode response: %+v", err)
	}

	var want Miner
	want = Miner{Name: "Борис", Energy: 2, Production: 3}

	if got != want {
		t.Errorf("got miner=%+v; want %+v", got, want)
	}

	if response.Code != http.StatusOK || contentType != "application/json" {
		t.Errorf("got code=%d and content type=%q; want code=200 and \"application/json\"", response.Code, contentType)
	}
}

func TestMinerHandlerNotFound(t *testing.T) {
	mine := Mine{Coal: 0, Reserves: 1, Closed: true}
	api := API{Mine: &mine}

	request := httptest.NewRequest(http.MethodGet, "/miner?id=999", nil)
	response := httptest.NewRecorder()
	api.minerHandler(response, request)

	if response.Code != http.StatusNotFound || response.Body.String() != "Шахтёр не найден\n" {
		t.Errorf("got code=%d and error=%q; want 404 and \"Шахтёр не найден\"", response.Code, response.Body.String())
	}
}

func TestMinerHandlerInvalidID(t *testing.T) {
	mine := Mine{Coal: 0, Reserves: 1, Closed: true}
	api := API{Mine: &mine}
	params := []string{"/miner", "/miner?id=abc"}
	for _, param := range params {
		request := httptest.NewRequest(http.MethodGet, param, nil)
		response := httptest.NewRecorder()
		api.minerHandler(response, request)

		if response.Code != http.StatusBadRequest || response.Body.String() != "Ошибка запроса\n" {
			t.Errorf("request=%q: got code=%d, body=%q; want 400 and body=%q", param, response.Code, response.Body.String(), "Ошибка запроса\n")
		}
	}
}

func TestMinerHandlerMethodNotAllowed(t *testing.T) {
	mine := Mine{Coal: 4, Reserves: 7, Closed: true}
	api := API{Mine: &mine}

	request := httptest.NewRequest(http.MethodPost, "/miner?id=101", nil)
	response := httptest.NewRecorder()
	api.minerHandler(response, request)

	if response.Code != http.StatusMethodNotAllowed || response.Header().Get("Allow") != http.MethodGet {
		t.Errorf("got response code=%d and response header=%q; want 405 and Allow=\"GET\"", response.Code, response.Header().Get("Allow"))
	}
}

func TestRouterMinerFound(t *testing.T) {
	mine := Mine{Coal: 4, Reserves: 7, Closed: true}
	mine.AddMiner(101, &Miner{Name: "Борис", Energy: 2, Production: 3})
	api := API{Mine: &mine}

	router := newRouter(api)

	request := httptest.NewRequest(http.MethodGet, "/miner?id=101", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	contentType := response.Header().Get("Content-Type")

	var got Miner
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&got); err != nil {
		t.Fatalf("decode response: %+v", err)
	}

	var want Miner
	want = Miner{Name: "Борис", Energy: 2, Production: 3}

	if got != want {
		t.Errorf("got miner=%+v; want %+v", got, want)
	}

	if response.Code != http.StatusOK || contentType != "application/json" {
		t.Errorf("got code=%d and content type=%q; want code=200 and \"application/json\"", response.Code, contentType)
	}
}
