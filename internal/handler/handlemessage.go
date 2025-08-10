package handle

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2/log"
)

func (h *Handler) HandlerImplMessage(message []byte, topic kafka.TopicPartition, cn int) error {
	log.Infof("Consumer #%d, Message from kafka with offset %d '%s' on partition %d", cn, topic.Offset, string(message), topic.Partition)
	return nil
}
