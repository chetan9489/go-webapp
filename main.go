package main

import (
	"net/http"
	"projects/go-webapp/server/api/user"

	log "github.com/sirupsen/logrus"
)

const (
	connStr = "user=postgres dbname=postgres host=localhost password=Che!@#123 sslmode=disable"
)

func main() {
	log.Info("Go ORM Tutorial")

	// register static files handle '/index.html -> client/index.html'
	http.Handle("/", http.FileServer(http.Dir("client")))
	// register RESTful endpoint handler for '/users/'
	http.Handle("/users/", &user.UserAPI{})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
