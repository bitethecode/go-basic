package server

import (
	"fmt"
	"net/http"

	"go-basic/internal/pkg/handler"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func Run() {
	fmt.Println("app starts")

	router := mux.NewRouter()
	// login & register with jwt
	router.HandleFunc("/api/v1/login", handler.Login).Methods("POST")
	router.HandleFunc("/api/v1/register", handler.Register).Methods("POST")

	// user
	router.HandleFunc("/api/v1/users", handler.GetUsers).Methods("GET")

	// books
	router.HandleFunc("/api/v1/books", handler.GetBooks).Methods("GET")
	router.HandleFunc("/api/v1/books", handler.PostBook).Methods("POST")
	http.ListenAndServe(":8080", router)
}
