package producers

import(
	"fmt"
	"github.com/IBM/sarama"
)

const KafkaServerAddress = "192.168.0.7:9092"

func SetupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{KafkaServerAddress}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup producer: %w", err)
	}
	return producer, nil
}