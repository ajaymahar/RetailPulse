package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ajaymahar/RetailPulse/internal/datastore"
	"github.com/ajaymahar/RetailPulse/internal/rest"
	"github.com/ajaymahar/RetailPulse/internal/service"
	"github.com/gorilla/mux"
)

func Start() {

	// just for now
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)
	//wiring the app
	// could be replaced with real DB
	repo := datastore.NewJobRepositoryStub()

	// Service port
	svc := service.NewDefaultJobService(repo)

	// mux router
	router := mux.NewRouter()

	// register all the routes
	rest.NewJobHandler(svc).Register(router)

	fmt.Println("Jobservice App started")
	s := &http.Server{
		Addr:              "localhost:8080",
		Handler:           router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       3 * time.Second,
		ErrorLog:          logger,
	}

	// NOTE:
	// Server to start listening requests
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			s.ErrorLog.Fatal(err.Error())
		}
	}()

	// NOTE:
	// chan to keep listening for os sign
	sigChan := make(chan os.Signal, 1)

	// NOTE:
	// any os action interrupt or kill will will notify the sigChan
	// signal.Notify(sigChan, os.Interrupt)
	// signal.Notify(sigChan, os.Kill)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// NOTE: waiting to any os signal interruption to do gracefull shoutdown
	sig := <-sigChan
	logger.Println("Recieved signal", sig)

	// NOTE:
	// context with timeout of 30 sec, to do gracefull shoutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// TODO:
	// gracefull Shutdown
	err := s.Shutdown(ctx)
	if err != nil {
		s.ErrorLog.Println("Error while shoutting down the server", err.Error())
	}
}
