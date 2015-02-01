package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/monban/timecard/store"
	"log"
	"net/http"
	"os"
)

func main() {
	store := board.Store{}
	store.Init()

	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
	}
	err := handler.SetRoutes(
		&rest.Route{"GET", "/employees", store.GetAllEmployees},
		&rest.Route{"POST", "/employees", store.PostEmployee},

		&rest.Route{"GET", "/transastions", store.GetAllTransactions},
		&rest.Route{"POST", "/transactions", store.PostTransaction},

		&rest.Route{"GET", "/locations", store.GetAllLocations},
		&rest.Route{"POST", "/locations", store.PostLocation},
	)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/api/", http.StripPrefix("/api", &handler))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", func(output http.ResponseWriter, request *http.Request) {
		http.ServeFile(output, request, "./static/board.html")
	})
	listenAddr := "0.0.0.0:" + get_port()
	log.Print("Starting server on " + listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func get_port() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	} else {
		log.Print("PORT variable not set, defaulting to 8080")
		return "8080"
	}
}
