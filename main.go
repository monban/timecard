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
		&rest.Route{"GET", "/transastions", store.GetAllTransactions},
		&rest.Route{"POST", "/employees", store.PostEmployee},
		&rest.Route{"POST", "/transactions", store.PostTransaction},
	)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/api/", http.StripPrefix("/api", &handler))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", func(output http.ResponseWriter, request *http.Request) {
		http.ServeFile(output, request, "./static/board.html")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Employee struct {
	Id          int64 `json:"id"`
	Name        string
	Location_id int64
}

type Transaction struct {
	Id          int64 `json:"id"`
	Employee_id int64
}

type Location struct {
	Id      int64 `json:"id"`
	Name    string
	OnClock bool
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
	store.DB.AutoMigrate(&Transaction{})
	store.DB.AutoMigrate(&Location{})
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

func (store *Store) GetAllTransactions(output rest.ResponseWriter, request *rest.Request) {
	transactions := []Transaction{}
	store.DB.Find(&transactions)
	output.WriteJson(&transactions)
}

func (store *Store) PostTransaction(output rest.ResponseWriter, request *rest.Request) {
	transaction := Transaction{}
	err := request.DecodeJsonPayload(&transaction)
	if err != nil {
		rest.Error(output, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := store.DB.Save(&transaction).Error; err != nil {
		rest.Error(output, err.Error(), http.StatusInternalServerError)
		return
	}
}
