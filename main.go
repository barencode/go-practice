package main

import (
	"fmt"
	"net/http"
)

func main() {
	mine := Mine{Reserves: 10}
	mine.AddMiner(101, &Miner{Name: "Борис", Energy: 2, Production: 3})
	api := &API{Mine: &mine}

	router := newRouter(api)

	fmt.Println("Сервер: http://localhost:8080")

	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		fmt.Println("Ошибка сервера:", err)
	}
}
