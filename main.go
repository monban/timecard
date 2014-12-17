package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
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
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Employee struct {
	Id           int64 `json:"id"`
	Name         string
	LocationId   int64
	Transactions []Transaction
}

func (employee *Employee) AfterCreate(db *gorm.DB) (err error) {
	location := Location{}
	db.First(&location)
	db.Model(employee).
		Association("Transaction").
		Append(Transaction{LocationId: location, CreatedAt: time.Now()})
	return
}

type EmployeePresenter struct {
	Id         int64 `json:"id"`
	Name       string
	Location   string
	ReturnTime time.Time
}

type Transaction struct {
	Id         int64 `json:"id"`
	EmployeeId int64
	CreatedAt  time.Time
	LocationId Location
	ReturnTime time.Time
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
	var count int
	store.DB.Model(Location{}).Count(&count)
	if count == 0 {
		store.DB.Create(Location{
			Name:    "In",
			OnClock: true,
		})
	}
}

func (store *Store) GetAllEmployees(output rest.ResponseWriter, request *rest.Request) {
	employees := []Employee{}
	store.DB.Find(&employees)
	presentedEmployees := make([]EmployeePresenter, len(employees))
	for i, employee := range employees {
		transaction := Transaction{}
		store.DB.Where("employee_id = ?", employee.Id).Last(&transaction)
		location := Location{}
		store.DB.Model(&transaction).Related(&location)
		presentedEmployees[i] = EmployeePresenter{
			Name:       employee.Name,
			Location:   location.Name,
			ReturnTime: transaction.ReturnTime,
		}
	}
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
func (store *Store) GetAllLocations(output rest.ResponseWriter, request *rest.Request) {
	locations := []Location{}
	store.DB.Find(&locations)
	output.WriteJson(&locations)
}

func (store *Store) PostLocation(output rest.ResponseWriter, request *rest.Request) {
	location := Location{}
	err := request.DecodeJsonPayload(&location)
	if err != nil {
		rest.Error(output, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := store.DB.Save(&location).Error; err != nil {
		rest.Error(output, err.Error(), http.StatusInternalServerError)
		return
	}
}
