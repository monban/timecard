package board

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

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
	output.WriteJson(&presentedEmployees)
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

type Employee struct {
	Id           int64 `json:"id"`
	Name         string
	LocationId   int64
	Transactions []Transaction
}

func (employee *Employee) AfterCreate(db *gorm.DB) (err error) {
	location := Location{}
	db.First(&location)
	transaction := Transaction{EmployeeId: employee.Id, LocationId: location.Id}
	db.Create(&transaction)
	return err
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
	LocationId int64
	ReturnTime time.Time
}

type Location struct {
	Id      int64 `json:"id"`
	Name    string
	OnClock bool
}
