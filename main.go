package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	store := Store{}
	store.Init()

	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
	}
	err := handler.SetRoutes(
		&rest.Route{"GET", "/employees", store.GetAllEmployees},
		&rest.Route{"POST", "/employees", store.PostEmployee},
	)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/api/", http.StripPrefix("/api", &handler))
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Employee struct {
	Id   int64 `json:"id"`
	Name string
}

type Store struct {
	DB gorm.DB
}

func (store *Store) Init() {
	var err error
	store.DB, err = gorm.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	store.DB.AutoMigrate(&Employee{})
}

func (store *Store) GetAllEmployees(output rest.ResponseWriter, request *rest.Request) {
	employees := []Employee{}
	store.DB.Find(&employees)
	output.WriteJson(&employees)
}

func (store *Store) PostEmployee(output rest.ResponseWriter, request *rest.Request) {
	employee := Employee{}
	err := request.DecodeJsonPayload(&employee)
	if err != nil {
		rest.Error(output, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := store.DB.Save(&employee).Error; err != nil {
		rest.Error(output, err.Error(), http.StatusInternalServerError)
		return
	}
}
