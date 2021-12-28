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

type Book struct {
	Id       string `json:"id"`
	Title    string `json:"title" validate:"required"`
	Subtitle string `json:"subtitle" validate:"required"`
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	db := pkgdb.OpenDb()

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Fatal("error", err)
	}

	var books []Book
	err = scan.Rows(&books, rows)
	if err != nil {
		log.Fatal("error", err)
	}
	w.Header().Set("Content-Type", "application/json")

	msg, err := json.Marshal(books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occured internally"))
		return
	}
	w.Write(msg)

	defer rows.Close()
	defer db.Close()
}

func PostBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	v := validator.New()

	if err = v.Struct(book); err != nil {
		errors := make(map[string][]string)

		for _, err := range err.(validator.ValidationErrors) {
			name := strings.ToLower(err.Field())
			switch err.Tag() {
			case "required":
				errors[name] = append(errors[name], "The "+name+" is required")
				break
			default:
				errors[name] = append(errors[name], "The "+name+" is invalid")
				break
			}
		}

		msg, err := json.Marshal(errors)
		if err != nil {
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

	sql := `INSERT INTO books (title, subtitle) VALUES ($1, $2)`
	_, err = db.Exec(sql, book.Title, book.Subtitle)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}
