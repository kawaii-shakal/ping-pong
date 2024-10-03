package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func startPongServer() {
	http.HandleFunc("/pong", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	log.Println("Pong сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startPingClient(pongAddress string, interval time.Duration) {
	for {
		resp, err := http.Get(pongAddress + "/pong")
		if err != nil {
			log.Printf("Ошибка при запросе к %s: %v\n", pongAddress, err)
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			log.Printf("Ответ от сервера: %s\n", string(body))
			resp.Body.Close()
		}
		time.Sleep(interval)
	}
}

func main() {
	modeFlag := flag.String("mode", "", "Режим работы: ping или pong")
	pongAddressFlag := flag.String("address", "http://localhost:8080", "Адрес сервера Pong для режима ping")
	intervalFlag := flag.Int("interval", 5, "Интервал в секундах между запросами в режиме ping")

	flag.Parse()

	mode := os.Getenv("MODE")
	if mode == "" {
		mode = *modeFlag
	}

	pongAddress := os.Getenv("PING_ADDRESS")
	if pongAddress == "" {
		pongAddress = *pongAddressFlag
	}

	intervalEnv := os.Getenv("PING_INTERVAL")
	interval := *intervalFlag
	if intervalEnv != "" {
		if parsedInterval, err := strconv.Atoi(intervalEnv); err == nil {
			interval = parsedInterval
		} else {
			log.Printf("Ошибка при преобразовании PING_INTERVAL: %v\n", err)
		}
	}

	if mode == "pong" {
		startPongServer()
	} else if mode == "ping" {
		log.Printf("Ping клиент запущен. Адрес сервера: %s, интервал: %d секунд\n", pongAddress, interval)
		startPingClient(pongAddress, time.Duration(interval)*time.Second)
	} else {
		fmt.Println("Укажите режим работы с помощью флага -mode или переменной среды MODE: ping или pong")
		os.Exit(1)
	}
}
