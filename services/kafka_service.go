package services

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/jinzhu/gorm"
	"github.com/suleiman/Personal-Budget-Manager/models"
)

// Function to produce notifications to Kafka
func produceNotificationsForUser(userID uint, db *gorm.DB, producer sarama.AsyncProducer) {
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
		producer.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(notificationBytes),
		}
	}
}
