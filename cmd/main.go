package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/itsyaboikris/go_document_store/api"
	"github.com/itsyaboikris/go_document_store/store"
)

func main() {
	ds := store.NewStore()
	router := mux.NewRouter()
	api.RegisterRoutes(router, ds)

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		log.Printf("Route: %s, Methods: %v\n", path, methods)
		return nil
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server starting on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
