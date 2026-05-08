package main

import (
	"fmt"
	"log"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func main() {
	http.HandleFunc("/ping", pingHandler)

	port := ":8080"
	fmt.Printf("Сервер запускається на порту%s...\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Помилка при запуску сервера: %v", err)
	}
}