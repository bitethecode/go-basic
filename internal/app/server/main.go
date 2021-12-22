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
	router.HandleFunc("/api/v1/users", handler.GetUsers).Methods("GET")
	router.HandleFunc("/api/v1/users", handler.PostUser).Methods("POST")

	http.ListenAndServe(":8080", router)
}
