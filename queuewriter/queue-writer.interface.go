package queuewriter

import "github.com/MichalMitros/feed-parser/models"

type QueueWriterInterface interface {
	WriteToQueue(queueName string, shopItems chan models.ShopItem) error
}
