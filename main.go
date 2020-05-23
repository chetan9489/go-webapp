package main

import (
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"projects/go-webapp/server/api/user"

	"github.com/Shopify/sarama"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const (
	connStr = "user=postgres dbname=postgres host=localhost password=Che!@#123 sslmode=disable"
)

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

func main() {
	log.Info("Go ORM Tutorial")

	// register static files handle '/index.html -> client/index.html'
	http.Handle("/", http.FileServer(http.Dir("client")))
	initialMigration()
	go recieveUserFromKafka()
	// register RESTful endpoint handler for '/users/'
	http.Handle("/users/", &user.UserAPI{})
	log.Fatal(http.ListenAndServe(":8080", nil))
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

	topic := "newtest11"
	// How to decide partition, is it fixed value...?
	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Error("Error while consuming the topic")
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case err := <-consumer.Errors():
			log.Error(err)
			continue
		case msg := <-consumer.Messages():
			log.Info(string(msg.Value))
			saveUserToDB(string(msg.Value))
		case <-signals:
			log.Error("Interrupt is detected")
			break
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
