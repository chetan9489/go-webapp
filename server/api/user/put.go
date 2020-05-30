package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// updateUser godoc
// @Summary Updates user based on given ID
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {object} models.User
// @Router /users/{id} [put]
func updateUser(w http.ResponseWriter, r *http.Request) {
	log.Info("Update User Endpoint Hit")

	w.Header().Set("Content-Type", "application/json")
	jd := json.NewDecoder(r.Body)

	aUser := &User{}
	err := jd.Decode(aUser)
	if nil != err {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	name := aUser.Name
	email := aUser.Email
	password := aUser.Password

	var user User
	err = db.Where("Email = ?", email).Find(&user).Error
	if err != nil {
		log.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	log.Info("ID: ", user.ID)
	log.Info("Name: ", user.Name)
	log.Info("Password: ", user.Password)

	user.Name = name
	user.Password = password

	db.Save(&user)
	log.Info("Successfully Updated User")
}
