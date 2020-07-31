package process_service

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type ClientIdData struct {
	Count     int
	Timestamp time.Time
}

const durationThreshold = 5
const requestsThreshold = 5

type ProcessService struct {
	Cache map[int]ClientIdData
}

var processServiceInstance *ProcessService
var processServiceOnce sync.Once

func NewProcessService() *ProcessService {
	processServiceOnce.Do(func() {
		processServiceInstance = &ProcessService{
			Cache: make(map[int]ClientIdData),
		}
	})
	return processServiceInstance
}

func (p ProcessService) ProcessClients(clientId int, ch chan bool) {
	msg := fmt.Sprintf("Process request for the client ID %d", clientId)
	log.Println(msg)
	if data, isExist := p.Cache[clientId]; isExist {
		duration := time.Now().Sub(data.Timestamp)
		secs := duration.Seconds()
		if secs < durationThreshold && p.Cache[clientId].Count < requestsThreshold {
			data.Count++
			p.Cache[clientId] = data
			msg = fmt.Sprintf("Client ID %d: return true 1. Seconds since last request %f. Count %d",
				clientId, secs, p.Cache[clientId].Count)
			log.Println(msg)
			ch <- true
		} else if secs >= durationThreshold {
			p.Cache[clientId] = ClientIdData{Count: 1, Timestamp: time.Now()}
			msg = fmt.Sprintf("Client ID %d: return true 2. Seconds since last request %f. Count %d",
				clientId, secs, p.Cache[clientId].Count)
			log.Println(msg)
			ch <- true
		} else {
			msg = fmt.Sprintf("Client ID %d: return false. Seconds since last request %f. Count %d",
				clientId, secs, data.Count)
			log.Println(msg)
			ch <- false
		}
	} else {
		p.Cache[clientId] = ClientIdData{Count: 1, Timestamp: time.Now()}
		msg = fmt.Sprintf("Client ID %d: return true 3. New client ID", clientId)
		log.Println(msg)
		ch <- true
	}
}
