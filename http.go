package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type API struct {
	Mine *Mine
}

type MineStatus struct {
	Coal     int  `json:"coal"`
	Reserves int  `json:"reserves"`
	Closed   bool `json:"closed"`
}

func newRouter(api API) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/health", healthHandler)
	router.HandleFunc("/miner", api.minerHandler)
	router.HandleFunc("/status", api.statusHandler)

	return router
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Шахта на связи")
}

func (api API) statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	status := MineStatus{
		Coal:     api.Mine.Coal,
		Reserves: api.Mine.Reserves,
		Closed:   api.Mine.Closed,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (api API) minerHandler(w http.ResponseWriter, r *http.Request) {
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

	miner, err := api.Mine.FindMiner(id)
	if err != nil {
		http.Error(w, "Шахтёр не найден", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(miner)
}
