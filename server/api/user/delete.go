package user

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// doDelete deletes a user from the db using the path '/users/id', eg: /users/2
func deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Info("Delete User Endpoint Hit")

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	vars := mux.Vars(r)
	email := vars["email"]

	var user User
	db.Where("Email = ?", email).Find(&user)
	db.Delete(&user)

	log.Info("Successfully Deleted User")
}
