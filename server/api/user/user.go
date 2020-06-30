package user

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

const (
	connStr = "user=postgres dbname=postgres host=localhost password=Che!@#123 sslmode=disable"
)

type UserAPI struct{}
type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"not null"`
	Password string `gorm:"not null"`
}

func (u *UserAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))

	switch r.Method {
	case http.MethodGet:
		allUsers(w, r)
	case http.MethodPost:
		newUser(w, r)
	case http.MethodDelete:
		deleteUser(w, r)
	case http.MethodPut:
		updateUser(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unsupported method '%v' to %v\n", r.Method, r.URL)
		log.Printf("Unsupported method '%v' to %v\n", r.Method, r.URL)
	}
}
