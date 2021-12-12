package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ajaymahar/RetailPulse/internal/datastore"
	"github.com/ajaymahar/RetailPulse/internal/rest"
	"github.com/ajaymahar/RetailPulse/internal/service"
	"github.com/gorilla/mux"
)

func Start() {

	//wiring the app
	// could be replaced with real DB
	repo := datastore.NewJobRepositoryStub()

	// Service port
	svc := service.NewDefaultJobService(repo)

	// mux router
	router := mux.NewRouter()

	// register all the routes
	rest.NewJobHandler(svc).Register(router)

	//TODO: implement custome server
	// custom server with timeout params and shout down gracefully
	fmt.Println("App started")
	log.Fatal(http.ListenAndServe(":8080", router))
}
