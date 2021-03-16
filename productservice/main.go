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

type Product struct {
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
		Path("/products/{id}").
		HandlerFunc(app.getProductById)

	app.Router.
		Methods("GET").
		Path("/products").
		HandlerFunc(app.getProducts)

	app.Router.
		Methods("POST").
		Path("/products").
		HandlerFunc(app.addProduct)
}

func (app *App) getProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Fatal("No ID in the path")
	}

	product := &Product{}
	err := app.Database.QueryRow("SELECT * FROM products WHERE id = ?", id).Scan(&product.id, &product.name, &product.price)
	if err != nil {
		log.Fatal("Database SELECT failed")
	}

	log.Println("You fetched a thing!")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		panic(err)
	}
}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	res, err := app.Database.Query("SELECT * FROM products")
	if err != nil {
		log.Fatal(err)
	}

	var products []Product

	for res.Next() {
		product := Product{}
		err := res.Scan(&product.id, &product.name, &product.price)

		if err != nil {
            log.Fatal(err)
        }

		fmt.Println(product)

		products = append(products, product)
	}

	if err := res.Err(); err != nil {
        fmt.Println(err)
        return
    }

	jData, err := json.Marshal(products)

	if err != nil {
		// handle error
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jData)

	fmt.Println(products)

    return
}

func (app *App) addProduct(w http.ResponseWriter, r *http.Request) {
	var product *Product
  	_ = json.NewDecoder(r.Body).Decode(&product)
	fmt.Println(product)
	// _, err := app.Database.Exec("INSERT INTO `test` (name) VALUES ('myname')")
	// if err != nil {
	// 	log.Fatal("Database INSERT failed")
	// }

	log.Println("Add product called!")
	w.WriteHeader(http.StatusOK)
}

func CreateDatabase() (*sql.DB, error) {
	serverName := "127.0.0.1:3306"
	user := "root"
	password := "root"
	dbName := "products_db"

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

	log.Fatal(http.ListenAndServe(":8080", app.Router))
}