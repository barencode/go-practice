package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/barencode/go-practice/internal/httpapi"
	"github.com/barencode/go-practice/internal/postgres"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Println("Не задан DATABASE_URL")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	pool, err := postgres.Open(ctx, databaseURL)
	cancel()
	if err != nil {
		fmt.Println("Ошибка подключения к PostgreSQL:", err)
		return
	}
	defer pool.Close()

	fmt.Println("PostgreSQL подключён")

	loadCtx, loadCancel := context.WithTimeout(context.Background(), time.Second*5)
	mine, err := postgres.LoadMine(loadCtx, pool, 1)
	loadCancel()
	if err != nil {
		fmt.Println("Ошибка загрузки шахты:", err)
		return
	}

	api, err := httpapi.NewAPI(mine)
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
