package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Payment struct {
	id int 		`json:"id"`
	name string `json:"name"`
	price int 	`json:"price"`
}

type App struct {
	Router   *mux.Router
	Database *sql.DB
}

func (app *App) SetupRouter() {
	app.Router.
		Methods("GET").
		Path("/payments/{id}").
		HandlerFunc(app.getPaymentById)

	app.Router.
		Methods("GET").
		Path("/payments").
		HandlerFunc(app.getPayments)

	app.Router.
		Methods("POST").
		Path("/payments").
		HandlerFunc(app.addPayment)
}

func (app *App) getPaymentById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Fatal("No ID in the path")
	}

	payment := &Payment{}
	err := app.Database.QueryRow("SELECT * FROM payments WHERE id = ?", id).Scan(&payment.id, &payment.name, &payment.price)
	if err != nil {
		log.Fatal("Database SELECT failed")
	}

	log.Println("You fetched a thing!")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(payment); err != nil {
		panic(err)
	}
}

func (app *App) getPayments(w http.ResponseWriter, r *http.Request) {
	res, err := app.Database.Query("SELECT * FROM payments")
	if err != nil {
		log.Fatal(err)
	}

	var payments []Payment

	for res.Next() {
		payment := Payment{}
		err := res.Scan(&payment.id, &payment.name, &payment.price)

		if err != nil {
            log.Fatal(err)
        }

		fmt.Println(payment)

		payments = append(payments, payment)
	}

	if err := res.Err(); err != nil {
        fmt.Println(err)
        return
    }

	jData, err := json.Marshal(payments)

	if err != nil {
		// handle error
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jData)

	fmt.Println(payments)

    return
}

func (app *App) addPayment(w http.ResponseWriter, r *http.Request) {
	var payment *Payment
  	_ = json.NewDecoder(r.Body).Decode(&payment)
	fmt.Println(payment)

	log.Println("Add payment called!")
	w.WriteHeader(http.StatusOK)
}

func CreateDatabase() (*sql.DB, error) {
	serverName := "127.0.0.1:3306"
	user := "root"
	password := "root"
	dbName := "payments_db"

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true", user, password, serverName, dbName)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	database, err := CreateDatabase()
	if err != nil {
		log.Fatal("Database connection failed: %s", err.Error())
	}

	app := &App{
		Router:   mux.NewRouter().StrictSlash(true),
		Database: database,
	}

	app.SetupRouter()

	log.Fatal(http.ListenAndServe(":8081", app.Router))
}