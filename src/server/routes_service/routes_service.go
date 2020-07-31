package routes_service

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/palantir/stacktrace"
	"io"
	"log"
	"net/http"
	"server/handlers"
	"sync"
	"time"
)

func InitRoutes() error {
	log.Println("Configuring routes")
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/").Subrouter()
	subRouter.HandleFunc("/", handlers.ProcessClientHandler)
	http.Handle("/", router)

	log.Println("Server is listening in the port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to init the routes")
	}

	log.Printf("main: starting HTTP server")
	httpServerExitDone := &sync.WaitGroup{}

	httpServerExitDone.Add(1)
	srv := startHttpServer(httpServerExitDone)

	log.Printf("main: serving for 10 seconds")

	time.Sleep(10 * time.Second)

	log.Printf("main: stopping HTTP server")

	// now close the server gracefully ("shutdown")
	// timeout could be given with a proper context
	// (in real world you shouldn't use TODO()).
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}

	// wait for goroutine started in startHttpServer() to stop
	httpServerExitDone.Wait()

	log.Printf("main: done. exiting")

	return nil
}

func startHttpServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: ":8080"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world\n")
	})

	go func() {
		defer wg.Done() // let main know we are done cleaning up

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}

/*
func InitRoutes() error {
	log.Println("Configuring routes")
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/").Subrouter()
	subRouter.HandleFunc("/", handlers.ProcessClientHandler)
	http.Handle("/", router)

	log.Println("Server is listening in the port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to init the routes")
	}
	return nil
}
*/
