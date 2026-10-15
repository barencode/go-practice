package main

import (
	"fmt"
	"net/http"
)

func main() {
	mine := Mine{Reserves: 10}
	mine.AddMiner(101, &Miner{Name: "Борис", Energy: 2, Production: 3})
	api := API{Mine: &mine}

	http.HandleFunc("/status", api.statusHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/miner", api.minerHandler)
	fmt.Println("Сервер: http://localhost:8080")

	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		fmt.Println("Ошибка сервера:", err)
	}
}
