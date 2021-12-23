package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	pkgdb "go-basic/internal/pkg/db"

	"github.com/blockloop/scan"
	"github.com/go-playground/validator/v10"
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
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

	msg, err := json.Marshal(users); if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occured internally"))
		return
	}
	w.Write(msg)

	defer rows.Close()
	defer db.Close()
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	errors := make(map[string]string)
	v := validator.New()

	if err = v.Struct(user); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors[strings.ToLower(err.Field())] = fmt.Sprintf("%s is %s %s", err.Field(), err.Tag(), err.Param())
		}

		msg, err := json.Marshal(errors); if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("An error occured internally"))
			return

		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}

	db := pkgdb.OpenDb()

	sql := `INSERT INTO users (username, password, email) VALUES ($1, $2, $3)`
	_, err = db.Exec(sql, user.Username, user.Password, user.Email)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}
