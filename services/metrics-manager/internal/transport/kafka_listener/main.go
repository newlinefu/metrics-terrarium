package kafka_listener

import (
	"github.com/IBM/sarama"
	"log"
	"metricsTerrarium/lib"
)

type Kafka struct {
	Consumer sarama.PartitionConsumer
}

type KafkaProperties struct {
	Config *lib.Config
}

func (k Kafka) Start(properties KafkaProperties) Kafka {
	consumer, err := sarama.NewConsumer([]string{properties.Config.KafkaAddress}, nil)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	} else {
		log.Printf("Consumer created succesfully")
	}
	defer consumer.Close()
	partConsumer, err := consumer.ConsumePartition("metrics", 0, sarama.OffsetNewest)
	k.Consumer = partConsumer
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumer.Close()

	return k
}
