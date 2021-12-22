package handler

import (
	"encoding/json"
	"log"
	"net/http"

	pkgdb "go-basic/internal/pkg/db"

	"github.com/blockloop/scan"
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := pkgdb.OpenDb()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal("error", err)
	}

	var users []User
	err = scan.Rows(&users, rows)
	if err != nil {
		log.Fatal("error", err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)

	defer rows.Close()
	defer db.Close()
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	db := pkgdb.OpenDb()

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sql := `INSERT INTO users (username, password, email) VALUES ($1, $2, $3)`
	_, err = db.Exec(sql, user.Username, user.Password, user.Email)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}
