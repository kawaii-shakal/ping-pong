package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Структура для отправки имени хоста в POST-запросе
type PingPayload struct {
    Hostname string `json:"hostname"`
}

// Структура для ответа сервера pong
type PongResponse struct {
    ClientIP string `json:"client_ip"`
    Hostname string `json:"hostname"`
}

func startPongServer(pingPort string) {
    http.HandleFunc("/pong", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            // Получаем IP-адрес клиента
            clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
            if err != nil {
                log.Printf("Ошибка при получении IP-адреса клиента: %v\n", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            // Читаем тело запроса
            var payload PingPayload
            body, err := ioutil.ReadAll(r.Body)
            if err != nil {
                log.Printf("Ошибка при чтении тела запроса: %v\n", err)
                http.Error(w, "Bad Request", http.StatusBadRequest)
                return
            }
            if err := json.Unmarshal(body, &payload); err != nil {
                log.Printf("Ошибка при разборе JSON: %v\n", err)
                http.Error(w, "Bad Request", http.StatusBadRequest)
                return
            }

            // Формируем и возвращаем ответ
            response := PongResponse{
                ClientIP: clientIP,
                Hostname: payload.Hostname,
            }
            responseBody, _ := json.Marshal(response)
            w.Header().Set("Content-Type", "application/json")
            w.Write(responseBody)
			log.Printf("Ответ: %s\n", string(responseBody))
        } else {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    })

    log.Println("Pong сервер запущен на порту"+pingPort)
    log.Fatal(http.ListenAndServe(":"+pingPort, nil))
}

// Функция для режима "ping"
func startPingClient(pongAddress string, interval time.Duration, hostname string) {
    for {
        // Формируем тело запроса с именем хоста
        payload := PingPayload{Hostname: hostname}
        payloadBytes, err := json.Marshal(payload)
        if err != nil {
            log.Printf("Ошибка при сериализации JSON: %v\n", err)
            continue
        }

        // Отправляем POST-запрос
        resp, err := http.Post(pongAddress+"/pong", "application/json", bytes.NewBuffer(payloadBytes))
        if err != nil {
            log.Printf("Ошибка при запросе к %s: %v\n", pongAddress, err)
        } else {
            // Читаем ответ
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

	hostname := os.Getenv("HOSTNAME")
    if hostname == "" {
        log.Println("Переменная среды HOSTNAME не найдена, используем 'unknown'")
        hostname = "unknown"
    }

	pingPort := os.Getenv("PONG_PORT")
    if hostname == "" {
        log.Println("Переменная среды HOSTNAME не найдена, используем 'unknown'")
        hostname = "unknown"
    }

	if mode == "pong" {
		startPongServer(pingPort)
	} else if mode == "ping" {
		log.Printf("Ping клиент запущен. Адрес сервера: %s, интервал: %d секунд\n", pongAddress, interval, hostname)
		startPingClient(pongAddress, time.Duration(interval)*time.Second, hostname)
	} else {
		fmt.Println("Укажите режим работы с помощью флага -mode или переменной среды MODE: ping или pong")
		os.Exit(1)
	}
}
