package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/suleiman-oss/Personal-Budget-Manager/models"
	"gopkg.in/Shopify/sarama.v1"
)

// Function to produce notifications to Kafka
func produceNotificationsForUser(userID uint, db *gorm.DB, producer sarama.SyncProducer) {
	var notifications []models.Notification
	if err := db.Where("user_id = ?", userID).Find(&notifications).Error; err != nil {
		// Handle error
	}

	for _, notification := range notifications {
		// Convert notification to JSON bytes
		notificationBytes, err := json.Marshal(notification)
		if err != nil {
			// Handle error
		}

		// Send notification to Kafka topic based on userID
		topic := fmt.Sprintf("%d", userID)
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(notificationBytes),
		}
		_, _, err = producer.SendMessage(msg)
		if err != nil {
			log.Fatal("Error sending message")
		}
		time.Sleep(time.Second)
	}
}
