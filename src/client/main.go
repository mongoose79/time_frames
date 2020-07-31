package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

const logFile = "TimeFrames_client.log"

func main() {
	initLog()
	log.Println("Starting time frames client...")

	log.Println("Enter the number of clients")
	var numOfClients int
	_, err := fmt.Scan(&numOfClients)
	if err != nil {
		log.Println("Wrong number of clients were entered")
		return
	}

	var wg sync.WaitGroup
	closeAllGoroutines := make(chan struct{})
	for i := 1; i <= numOfClients; i++ {
		wg.Add(1)
		go processClient(i, closeAllGoroutines, &wg)
	}

	var tmp int
	_, err = fmt.Scan(&tmp)
	if err != nil {
		log.Println("Error occurred")
		return
	}
	close(closeAllGoroutines)
	wg.Wait()
	log.Println("Time frames client was completed successfully")
}

func processClient(clientId int, closeGoroutine chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	shouldProcess := true
	for shouldProcess {
		select {
		case <-closeGoroutine:
			msg := fmt.Sprintf("Goroutine with channel ID %d will be closed", clientId)
			log.Println(msg)
			shouldProcess = false
		default:
			url := fmt.Sprintf("http://localhost:8080/?clientId=%d", clientId)
			response, err := http.Get(url)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to process request for client ID %d", clientId)
				log.Println(errMsg)
				return
			}
			msg := fmt.Sprintf("Process request for client ID %d. Status: %d", clientId, response.StatusCode)
			log.Println(msg)

			msToWait := random(300, 800)
			msg = fmt.Sprintf("Sleep for %d milliseconds", msToWait)
			log.Println(msg)
			time.Sleep(time.Duration(msToWait) * time.Millisecond)
		}
	}
	msg := fmt.Sprintf("Close goroutine with channel ID %d", clientId)
	log.Println(msg)
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func initLog() {
	fmt.Println("Start initializing the log")
	logFile, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Failed to create log file")
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
