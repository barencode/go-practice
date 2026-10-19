package main

import (
	"fmt"
	"net/http"

	"github.com/barencode/go-practice/internal/httpapi"
	"github.com/barencode/go-practice/internal/mining"
)

func main() {
	mine := mining.Mine{Reserves: 10}
	mine.AddMiner(101, &mining.Miner{Name: "Борис", Energy: 2, Production: 3})
	api, err := httpapi.NewAPI(&mine)
	if err != nil {
		fmt.Println("Ошибка создания API:", err)
		return
	}

	router := httpapi.NewRouter(api)

	fmt.Println("Сервер: http://localhost:8080")

	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		fmt.Println("Ошибка сервера:", err)
	}
}
