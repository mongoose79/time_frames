package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"server/routes_service"
)

const logFile = "TimeFrames_server.log"

func main() {
	initLog()
	log.Println("Starting time frames server...")
	err := routes_service.InitRoutes()
	if err != nil {
		log.Println(fmt.Sprintf("Failed to init the routes: %v", err))
	}
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
