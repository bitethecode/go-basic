package handler

import (
	"encoding/json"
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

	v := validator.New()

	if err = v.Struct(user); err != nil {
		errors := make(map[string][]string)

		for _,err := range err.(validator.ValidationErrors){
            name := strings.ToLower(err.Field())
            switch err.Tag() {
            case "required":
                errors[name] = append(errors[name], "The "+name+" is required")
                break
            case "email":
                errors[name] = append(errors[name], "The "+name+" should be a valid email")
                break
            case "eqfield":
                errors[name] = append(errors[name], "The "+name+" should be equal to the "+err.Param())
                break
            default:
                errors[name] = append(errors[name], "The "+name+" is invalid")
                break
            }
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
