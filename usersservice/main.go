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

type User struct {
	id int 			`json:"id"`
	login string 	`json:"login"`
	password int 	`json:"password"`
}

type App struct {
	Router   *mux.Router
	Database *sql.DB
}

func (app *App) SetupRouter() {
	app.Router.
		Methods("GET").
		Path("/users/{id}").
		HandlerFunc(app.getUserById)

	app.Router.
		Methods("GET").
		Path("/users").
		HandlerFunc(app.getUsers)

	app.Router.
		Methods("POST").
		Path("/users").
		HandlerFunc(app.addUser)
}

func (app *App) getUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Fatal("No ID in the path")
	}

	user := &User{}
	err := app.Database.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user.id, &user.login, &user.password)
	if err != nil {
		log.Fatal("Database SELECT failed")
	}

	log.Println("You fetched a thing!")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

func (app *App) getUsers(w http.ResponseWriter, r *http.Request) {
	res, err := app.Database.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}

	var users []User

	for res.Next() {
		user := User{}
		err := res.Scan(&user.id, &user.login, &user.password)

		if err != nil {
            log.Fatal(err)
        }

		fmt.Println(user)

		users = append(users, user)
	}

	if err := res.Err(); err != nil {
        fmt.Println(err)
        return
    }

	jData, err := json.Marshal(users)

	if err != nil {
		// handle error
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jData)

	fmt.Println(users)

    return
}

func (app *App) addUser(w http.ResponseWriter, r *http.Request) {
	var user *User
  	_ = json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(user)
	// _, err := app.Database.Exec("INSERT INTO `test` (name) VALUES ('myname')")
	// if err != nil {
	// 	log.Fatal("Database INSERT failed")
	// }

	log.Println("Add user called!")
	w.WriteHeader(http.StatusOK)
}

func CreateDatabase() (*sql.DB, error) {
	serverName := "127.0.0.1:3306"
	user := "root"
	password := "root"
	dbName := "users_db"

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

	log.Fatal(http.ListenAndServe(":8082", app.Router))
}