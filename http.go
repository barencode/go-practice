package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type API struct {
	Mine *Mine
	mtx  sync.Mutex
}

type MineStatus struct {
	Coal     int  `json:"coal"`
	Reserves int  `json:"reserves"`
	Closed   bool `json:"closed"`
}

type MineRequest struct {
	MinerID int `json:"miner_id"`
}

type MineResponse struct {
	Coal int `json:"coal"`
}

type RestRequest struct {
	MinerID int `json:"miner_id"`
}

func newAPI(mine *Mine) (*API, error) {
	if mine == nil {
		return nil, errors.New("шахта не задана")
	}
	return &API{Mine: mine}, nil
}

func newRouter(api *API) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/health", healthHandler)
	router.HandleFunc("/mine", api.mineHandler)
	router.HandleFunc("/miner", api.minerHandler)
	router.HandleFunc("/rest", api.restHandler)
	router.HandleFunc("/status", api.statusHandler)

	return router
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Шахта на связи")
}

func (api *API) statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	api.mtx.Lock()
	defer api.mtx.Unlock()

	status := MineStatus{
		Coal:     api.Mine.Coal,
		Reserves: api.Mine.Reserves,
		Closed:   api.Mine.Closed,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (api *API) minerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	input := r.URL.Query().Get("id")
	id, err := parseMinerID(input)
	if err != nil {
		http.Error(w, "Ошибка запроса", http.StatusBadRequest)
		return
	}

	api.mtx.Lock()
	defer api.mtx.Unlock()

	miner, err := api.Mine.FindMiner(id)
	if err != nil {
		http.Error(w, "Шахтёр не найден", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(miner)
}

func (api *API) mineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var input MineRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Ошибка запроса", http.StatusBadRequest)
		return
	}

	api.mtx.Lock()
	defer api.mtx.Unlock()

	coal, err := api.Mine.MineByID(input.MinerID)
	if errors.Is(err, ErrMineClosed) {
		http.Error(w, "Шахта закрыта", http.StatusConflict)
		return
	} else if errors.Is(err, ErrMinerNotFound) {
		http.Error(w, "Шахтёр не найден", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MineResponse{Coal: coal})
}

func (api *API) restHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var input RestRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Ошибка запроса", http.StatusBadRequest)
		return
	}

	api.mtx.Lock()
	defer api.mtx.Unlock()

	err = api.Mine.RestMiner(input.MinerID)
	if errors.Is(err, ErrMinerNotFound) {
		http.Error(w, "Шахтёр не найден", http.StatusNotFound)
		return
	} else if errors.Is(err, ErrNotEnoughCoal) {
		http.Error(w, "Недостаточно угля", http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RestRequest{MinerID: input.MinerID})
}
