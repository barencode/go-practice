package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", healthHandler)
	fmt.Println("Сервер: http://localhost:8080")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		fmt.Println("Ошибка сервера:", err)
	}
}
