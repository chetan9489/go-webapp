package user

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// deleteUser godoc
// @Summary Deletes user from DB based on User ID.
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {}
// @Router /users/{id} [delete]
func deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Info("Delete User Endpoint Hit")

	query := r.URL.Query()
	ID := query.Get("ID")
	log.Info("Id : ", ID)

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var users []User
	json.NewEncoder(w).Encode(db.Where("ID = ?", ID).Find(&users))
	log.Info(users)
	db.Where("ID = ?", ID).Delete(&users)

	log.Info("Successfully Deleted User")
}
