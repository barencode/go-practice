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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Шахта на связи")
}

func (api API) statusHandler(w http.ResponseWriter, r *http.Request) {
	status := MineStatus{
		Coal:     api.Mine.Coal,
		Reserves: api.Mine.Reserves,
		Closed:   api.Mine.Closed,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
