package handlers

import (
	"fmt"
	"log"
	"net/http"
	"server/process_service"
	"server/utils"
	"strconv"
)

func ProcessClientHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Process Ermetic client request was received. Starting...")
	clientIdStr := r.URL.Query().Get("clientId")
	if clientIdStr == "" {
		errMsg := "Failed source argument 'clientId'"
		log.Println(errMsg)
		utils.WriteJSON(errMsg, w, http.StatusBadRequest)
		return
	}

	var clientId int
	var err error
	if clientId, err = strconv.Atoi(clientIdStr); err != nil {
		errMsg := fmt.Sprintf("Failed to parse %s to int", clientIdStr)
		log.Println(errMsg)
	}
	msg := fmt.Sprintf("Start checking request for client ID %d", clientId)
	log.Println(msg)
	ch := make(chan bool)

	ps := process_service.NewProcessService()
	go ps.ProcessClients(clientId, ch)
	retVal := <-ch
	msg = fmt.Sprintf("Checking request for client ID %d was completed successfully. Return value: %t", clientId, retVal)
	log.Println(msg)
	if retVal {
		utils.WriteJSON("StatusOK", w, http.StatusOK)
		return
	}
	utils.WriteJSON("StatusServiceUnavailable", w, http.StatusServiceUnavailable)
}
