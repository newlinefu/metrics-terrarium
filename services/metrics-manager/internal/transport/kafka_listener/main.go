package kafka_listener

import (
	"github.com/IBM/sarama"
	"log"
	"metricsTerrarium/lib"
)

type Kafka struct {
	Consumer     *sarama.Consumer
	PartConsumer *sarama.PartitionConsumer
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
	partConsumer, err := consumer.ConsumePartition("metrics", 0, sarama.OffsetNewest)

	k.Consumer = &consumer
	k.PartConsumer = &partConsumer

	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}

	return k
}
