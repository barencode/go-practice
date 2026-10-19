package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/barencode/go-practice/internal/mining"
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
	mine := mining.Mine{Coal: 4, Reserves: 7, Closed: true}
	api := &API{Mine: &mine}

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
	mine := mining.Mine{Coal: 4, Reserves: 7, Closed: true}
	api := &API{Mine: &mine}

	request := httptest.NewRequest(http.MethodPost, "/status", nil)
	response := httptest.NewRecorder()
	api.statusHandler(response, request)

	if response.Code != http.StatusMethodNotAllowed || response.Header().Get("Allow") != http.MethodGet {
		t.Errorf("got response code=%d and response header=%q; want 405 and Allow=\"GET\"", response.Code, response.Header().Get("Allow"))
	}
}

func TestMinerHandlerFound(t *testing.T) {
	mine := mining.Mine{Coal: 4, Reserves: 7, Closed: true}
	mine.AddMiner(101, &mining.Miner{Name: "Борис", Energy: 2, Production: 3})
	api := &API{Mine: &mine}

	request := httptest.NewRequest(http.MethodGet, "/miner?id=101", nil)
	response := httptest.NewRecorder()
	api.minerHandler(response, request)
	contentType := response.Header().Get("Content-Type")

	var got mining.Miner
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&got); err != nil {
		t.Fatalf("decode response: %+v", err)
	}

	var want mining.Miner
	want = mining.Miner{Name: "Борис", Energy: 2, Production: 3}

	if got != want {
		t.Errorf("got miner=%+v; want %+v", got, want)
	}

	if response.Code != http.StatusOK || contentType != "application/json" {
		t.Errorf("got code=%d and content type=%q; want code=200 and \"application/json\"", response.Code, contentType)
	}
}

func TestMinerHandlerNotFound(t *testing.T) {
	mine := mining.Mine{Coal: 0, Reserves: 1, Closed: true}
	api := &API{Mine: &mine}

	request := httptest.NewRequest(http.MethodGet, "/miner?id=999", nil)
	response := httptest.NewRecorder()
	api.minerHandler(response, request)

	if response.Code != http.StatusNotFound || response.Body.String() != "Шахтёр не найден\n" {
		t.Errorf("got code=%d and error=%q; want 404 and \"Шахтёр не найден\"", response.Code, response.Body.String())
	}
}

func TestMinerHandlerInvalidID(t *testing.T) {
	mine := mining.Mine{Coal: 0, Reserves: 1, Closed: true}
	api := &API{Mine: &mine}
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
	mine := mining.Mine{Coal: 4, Reserves: 7, Closed: true}
	api := &API{Mine: &mine}

	request := httptest.NewRequest(http.MethodPost, "/miner?id=101", nil)
	response := httptest.NewRecorder()
	api.minerHandler(response, request)

	if response.Code != http.StatusMethodNotAllowed || response.Header().Get("Allow") != http.MethodGet {
		t.Errorf("got response code=%d and response header=%q; want 405 and Allow=\"GET\"", response.Code, response.Header().Get("Allow"))
	}
}

func TestRouterMinerFound(t *testing.T) {
	mine := mining.Mine{Coal: 4, Reserves: 7, Closed: true}
	mine.AddMiner(101, &mining.Miner{Name: "Борис", Energy: 2, Production: 3})
	api := &API{Mine: &mine}
	router := newRouter(api)

	request := httptest.NewRequest(http.MethodGet, "/miner?id=101", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	contentType := response.Header().Get("Content-Type")

	var got mining.Miner
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&got); err != nil {
		t.Fatalf("decode response: %+v", err)
	}

	var want mining.Miner
	want = mining.Miner{Name: "Борис", Energy: 2, Production: 3}

	if got != want {
		t.Errorf("got miner=%+v; want %+v", got, want)
	}

	if response.Code != http.StatusOK || contentType != "application/json" {
		t.Errorf("got code=%d and content type=%q; want code=200 and \"application/json\"", response.Code, contentType)
	}
}

func TestMineRouterSuccess(t *testing.T) {
	mine := mining.Mine{Coal: 4, Reserves: 5}
	mine.AddMiner(101, &mining.Miner{Name: "Борис", Energy: 1, Production: 3})
	api := &API{Mine: &mine}
	router := newRouter(api)

	body := strings.NewReader(`{"miner_id":101}`)
	request := httptest.NewRequest(http.MethodPost, "/mine", body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	contentType := response.Header().Get("Content-Type")

	if response.Code != http.StatusOK || contentType != "application/json" {
		t.Errorf("got code=%d and content type=%q; want code=200 and \"application/json\"", response.Code, contentType)
	}

	var got MineResponse
	if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got.Coal != 3 {
		t.Errorf("got mined coal=%d; want 3", got.Coal)
	}

	if mine.Coal != 7 || mine.Reserves != 2 || mine.Miners[101].Energy != 0 {
		t.Errorf("got coal=%d, reserves=%d and 1st miner energy=%d; want 7, 2 and 0", mine.Coal, mine.Reserves, mine.Miners[101].Energy)
	}
}

func TestMineRouterClosed(t *testing.T) {
	mine := mining.Mine{Closed: true, Coal: 4, Reserves: 5}
	mine.AddMiner(101, &mining.Miner{Name: "Борис", Energy: 1, Production: 3})
	api := &API{Mine: &mine}
	router := newRouter(api)

	body := strings.NewReader(`{"miner_id":101}`)
	request := httptest.NewRequest(http.MethodPost, "/mine", body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Errorf("got code=%d; want code=409", response.Code)
	}

	if response.Body.String() != "Шахта закрыта\n" {
		t.Errorf("got body=%q; want body=%q", response.Body.String(), "Шахта закрыта\n")
	}

	if mine.Coal != 4 || mine.Reserves != 5 || mine.Miners[101].Energy != 1 {
		t.Errorf("got coal=%d, reserves=%d and 1st miner energy=%d; want 4, 5 and 1", mine.Coal, mine.Reserves, mine.Miners[101].Energy)
	}
}

func TestMineRouterInvalidJSON(t *testing.T) {
	mine := mining.Mine{Coal: 4, Reserves: 5}
	mine.AddMiner(101, &mining.Miner{Name: "Борис", Energy: 1, Production: 3})
	api := &API{Mine: &mine}
	router := newRouter(api)

	body := strings.NewReader(`{"miner_id":`)
	request := httptest.NewRequest(http.MethodPost, "/mine", body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("got code=%d; want code=400", response.Code)
	}

	if response.Body.String() != "Ошибка запроса\n" {
		t.Errorf("got body=%q; want body=%q", response.Body.String(), "Ошибка запроса\n")
	}

	if mine.Coal != 4 || mine.Reserves != 5 || mine.Miners[101].Energy != 1 {
		t.Errorf("got coal=%d, reserves=%d and 1st miner energy=%d; want 4, 5 and 1", mine.Coal, mine.Reserves, mine.Miners[101].Energy)
	}
}

func TestMineRouterNotFound(t *testing.T) {
	mine := mining.Mine{Coal: 4, Reserves: 5}
	mine.AddMiner(101, &mining.Miner{Name: "Борис", Energy: 1, Production: 3})
	api := &API{Mine: &mine}
	router := newRouter(api)

	body := strings.NewReader(`{"miner_id":999}`)
	request := httptest.NewRequest(http.MethodPost, "/mine", body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("got code=%d; want code=404", response.Code)
	}

	if response.Body.String() != "Шахтёр не найден\n" {
		t.Errorf("got body=%q; want body=%q", response.Body.String(), "Шахтёр не найден\n")
	}

	if mine.Coal != 4 || mine.Reserves != 5 || mine.Miners[101].Energy != 1 {
		t.Errorf("got coal=%d, reserves=%d and 1st miner energy=%d; want 4, 5 and 1", mine.Coal, mine.Reserves, mine.Miners[101].Energy)
	}
}

func TestMineRouterMethodNotAllowed(t *testing.T) {
	mine := mining.Mine{Coal: 4, Reserves: 5}
	mine.AddMiner(101, &mining.Miner{Name: "Борис", Energy: 1, Production: 3})
	api := &API{Mine: &mine}
	router := newRouter(api)

	body := strings.NewReader(`{"miner_id":101}`)
	request := httptest.NewRequest(http.MethodGet, "/mine", body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed || response.Header().Get("Allow") != http.MethodPost {
		t.Errorf("got response code=%d and response header=%q; want 405 and Allow=\"POST\"", response.Code, response.Header().Get("Allow"))
	}

	if mine.Coal != 4 || mine.Reserves != 5 || mine.Miners[101].Energy != 1 {
		t.Errorf("got coal=%d, reserves=%d and 1st miner energy=%d; want 4, 5 and 1", mine.Coal, mine.Reserves, mine.Miners[101].Energy)
	}
}

func TestRouterMineRestAndMine(t *testing.T) {
	mine := mining.Mine{Reserves: 10}
	mine.AddMiner(101, &mining.Miner{Name: "Борис", Energy: 2, Production: 3})
	api := &API{Mine: &mine}
	router := newRouter(api)

	for _, path := range []string{"/mine", "/rest", "/mine"} {
		body := strings.NewReader(`{"miner_id":101}`)
		request := httptest.NewRequest(http.MethodPost, path, body)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		if response.Code != http.StatusOK {
			t.Fatalf("POST %s: got code=%d, body=%q; want 200", path, response.Code, response.Body.String())
		}
	}

	if mine.Coal != 4 || mine.Reserves != 4 || mine.Miners[101].Energy != 1 {
		t.Errorf("got coal=%d, reserves=%d and 1st miner energy=%d; want 4, 4, and 1", mine.Coal, mine.Reserves, mine.Miners[101].Energy)
	}
}

func TestRestRouterNotEnoughCoal(t *testing.T) {
	mine := mining.Mine{Coal: 1, Reserves: 5}
	mine.AddMiner(101, &mining.Miner{Name: "Борис", Energy: 0, Production: 3})
	api := &API{Mine: &mine}
	router := newRouter(api)

	path := "/rest"
	body := strings.NewReader(`{"miner_id":101}`)
	request := httptest.NewRequest(http.MethodPost, path, body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("POST %s: got code=%d, body=%q; want 409", path, response.Code, response.Body.String())
	}

	if response.Body.String() != "Недостаточно угля\n" {
		t.Errorf("got body=%q; want %q",
			response.Body.String(), "Недостаточно угля\n")
	}

	if mine.Coal != 1 || mine.Reserves != 5 || mine.Miners[101].Energy != 0 {
		t.Errorf("got coal=%d, reserves=%d and 1st miner energy=%d; want 1, 5 and 0", mine.Coal, mine.Reserves, mine.Miners[101].Energy)
	}
}

func TestNewAPINil(t *testing.T) {
	api, err := newAPI(nil)

	if api != nil || err == nil {
		t.Errorf("got api=%v and error=%v; want nil and not-nil error", api, err)
	}
}

func TestNewAPISuccess(t *testing.T) {
	mine := &mining.Mine{}
	api, err := newAPI(mine)

	if api == nil || err != nil {
		t.Fatalf("got api=%v and error=%v; want non-nil and nil error", api, err)
	}

	if api.Mine != mine {
		t.Errorf("got mine pointer=%p; want %p", api.Mine, mine)
	}
}
