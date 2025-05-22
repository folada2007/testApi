package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	_ "testApi/docs"
	"testApi/internal/handler"
	"testApi/internal/loger"
	"testApi/internal/repository"
)

// @title TestAPI
// @version 1.0
// @description CRUD api
// @host localhost:8080
// @BasePath /

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
}

func main() {
	loger.InitLogger()
	connStr := os.Getenv("DB_CONN_STRING")

	r := mux.NewRouter()

	pool, err := repository.GetConnectionPool(connStr)
	if err != nil {
		log.Println("Error creating connection pool:", err)
		panic(err)
	}

	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/", handler.GetFilteredUsersOrDefault(pool)).Methods("GET")
	r.HandleFunc("/PostUserName", handler.PostUserNameHandler(pool)).Methods("POST")
	r.HandleFunc("/UpdateUser/{id}", handler.UpdateUserDataHandler(pool)).Methods("PUT")
	r.HandleFunc("/DeleteUser/{id}", handler.DeleteUserDataHandler(pool)).Methods("DELETE")

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
