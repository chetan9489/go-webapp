package user

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// allUsers godoc
// @Summary Retrieves all users saved in the DB.
// @Produce json
// @Success 200 {object} models.User
// @Router /users/ [get]
func allUsers(w http.ResponseWriter, r *http.Request) {
	log.Info("All Users Endpoint Hit")

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var users []User
	db.Find(&users)
	json.NewEncoder(w).Encode(users)
}
