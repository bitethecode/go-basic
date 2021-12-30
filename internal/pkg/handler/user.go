package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	pkgdb "go-basic/internal/pkg/db"

	"github.com/blockloop/scan"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type Claims struct {
	jwt.StandardClaims
}

const SECRET_KEY = "gosecretkey"

func GetHash(pwd []byte) string {
	fmt.Print("get to here?")
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

func GenerateJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{})
	tokenString, err := token.SignedString([]byte("secretkey"))

	if err != nil {
		log.Println("Error in JWT token generation.")
		return "", err
	}

	return tokenString, nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	json.NewDecoder(r.Body).Decode(&user)

	db := pkgdb.OpenDb()

	var dbUser User
	sql := `SELECT * FROM users WHERE email = $1`
	rows, err := db.Query(sql, user.Email)

	err = scan.Row(&dbUser, rows)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found."))
		return
	}

	userPass := []byte(user.Password)
	dbPass := []byte(dbUser.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

	if passErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Wrong Password!"))
		return
	}

	jwtToken, err := GenerateJWT()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error in JWT otken generation."))
		return
	}

	msg, err := json.Marshal(fmt.Sprintf("token: %s", jwtToken))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occured internally"))
		return
	}
	w.Write(msg)
	defer rows.Close()
	defer db.Close()
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	v := validator.New()

	if err = v.Struct(user); err != nil {
		errors := make(map[string][]string)

		for _, err := range err.(validator.ValidationErrors) {
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

	// hash the password before saving
	sql := `INSERT INTO users (username, password, email) VALUES ($1, $2, $3)`
	_, err = db.Exec(sql, user.Username, GetHash([]byte(user.Password)), user.Email)

	fmt.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trasactional process"))
		return
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := pkgdb.OpenDb()

	rows, err := db.Query("SELECT * FROM users;")
	if err != nil {
		log.Fatal("error", err)
	}

	var users []User
	err = scan.Rows(&users, rows)
	if err != nil {
		log.Fatal("error", err)
	}
	w.Header().Set("Content-Type", "application/json")

	msg, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occured internally"))
		return
	}
	w.Write(msg)

	defer rows.Close()
	defer db.Close()
}
