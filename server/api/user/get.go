package user

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

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
	/* u, _ := db.Model(&User{}).Find(&User{}).Rows()
	for u.Next() {
		user := new(User)
		db.ScanRows(u, user)
		fmt.Printf("Got: Name: %v, Username: %v\n", user.Name, user.Email)
	} */
}
