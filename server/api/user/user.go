package user

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"

	"encoding/json"

	"github.com/Shopify/sarama"
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

func initialMigration() {
	log.Info("Connecting to SQL Db...")
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&User{})
}

/* var db = []*User{}
var nextUserID uint64
var lock sync.Mutex */

func (u *UserAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	initialMigration()
	go recieveUserFromKafka()

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

func recieveUserFromKafka() {

	log.Info("Start recieving User from kafka")

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Specify brokers address. This is default one
	brokers := []string{"localhost:9092"}

	// Create new consumer
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Error("Error while creating consumer")
		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			log.Error("Error while closing the consumer")
			panic(err)
		}
	}()

	topic := "newtest6"
	// How to decide partition, is it fixed value...?
	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Error("Error while consuming the topic")
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Get signal for finish
	doneCh := make(chan struct{})
	for {
		select {
		case err := <-consumer.Errors():
			log.Error(err)
		case msg := <-consumer.Messages():
			log.Info(string(msg.Value))
			saveUserToDB(string(msg.Value))
		case <-signals:
			log.Error("Interrupt is detected")
			doneCh <- struct{}{}
		}
	}
}

func saveUserToDB(userString string) {

	log.Info("Starting to save user to DB")

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var aUser User
	b := []byte(userString)
	error := json.Unmarshal(b, &aUser)
	if error != nil {
		log.Error("Error while unmarshalling json to struct")
		panic(error)
	}

	// CRUD = create, retrieve, update, and delete
	if err := db.Model(&User{}).Create(&User{Name: aUser.Name, Email: aUser.Email, Password: aUser.Password}).Error; err != nil {
		log.Warn(err)
	}
	log.Info("New User Successfully Created")

}
