package user

import (
	"encoding/json"
	"net/http"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

// newUser godoc
// @Summary Saves user on Kafka.
// @Produce json
// @Param object path User true "User"
// @Success 200 {}
// @Router /users/{User object} [post]
func newUser(w http.ResponseWriter, r *http.Request) {
	log.Info("New User Endpoint Hit")

	w.Header().Set("Content-Type", "application/json")
	jd := json.NewDecoder(r.Body)

	aUser := &User{}
	err := jd.Decode(aUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	saveUserToKafka(*aUser)

}

func saveUserToKafka(User User) {

	log.Info("Start saving to kafka")

	jsonObject, err := json.Marshal(User)

	jsonString := string(jsonObject)
	log.Info("User saved on Kafka")
	log.Info(jsonString)

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5

	// brokers := []string{"192.168.59.103:9092"}
	brokers := []string{"localhost:9092"}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		// Should not reach here
		log.Error("Error while creating a Producer")
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			// Should not reach here
			log.Error("Error while closing a Producer")
			panic(err)
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "newtest11"
	for _, word := range []string{string(jsonString)} {
		_, _, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder([]byte(word)),
		})
		if err != nil {
			log.Error("Error while writing on kafka")
			panic(err)
		}
	}

	log.Info("Message is stored in topic :", topic)
}
