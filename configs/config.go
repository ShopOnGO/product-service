package configs

import (
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
)

type Config struct {
	Db DbConfig
	KafkaProduct 	KafkaConsumerConfig
	KafkaVariant 	KafkaConsumerConfig
	KafkaMedia 	 	KafkaConsumerConfig
	KafkaProducer  	KafkaProducerConfig
}

type DbConfig struct {
	Dsn string
}

type KafkaConsumerConfig struct {
	Brokers []string
	Topic   string
	GroupID string
	ClientID string
}

type KafkaProducerConfig struct {
	Brokers []string
	Topic   map[string]string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file, using default config", err.Error())
	}

	brokersRaw := os.Getenv("KAFKA_BROKERS")
	brokers := strings.Split(brokersRaw, ",")

	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		KafkaProduct: KafkaConsumerConfig{
			Brokers: brokers,
			Topic:   os.Getenv("KAFKA_PRODUCT_TOPIC"),
			GroupID: os.Getenv("KAFKA_PRODUCT_GROUP_ID"),
			ClientID: os.Getenv("KAFKA_PRODUCT_CLIENT_ID"),
		},
		KafkaVariant: KafkaConsumerConfig{
			Brokers: brokers,
			Topic:   os.Getenv("KAFKA_VARIANT_TOPIC"),
			GroupID: os.Getenv("KAFKA_VARIANT_GROUP_ID"),
			ClientID: os.Getenv("KAFKA_VARIANT_CLIENT_ID"),
		},
		KafkaMedia: KafkaConsumerConfig{
			Brokers: brokers,
			Topic:   os.Getenv("KAFKA_MEDIA_TOPIC"),
			GroupID: os.Getenv("KAFKA_MEDIA_GROUP_ID"),
			ClientID: os.Getenv("KAFKA_MEDIA_CLIENT_ID"),
		},
		KafkaProducer: KafkaProducerConfig{
			Brokers: brokers,
			Topic:   parseKafkaTopics(os.Getenv("KAFKA_PRODUCER_TOPIC")),
		},
	}
}

func parseKafkaTopics(s string) map[string]string {
	topics := map[string]string{}
	pairs := strings.Split(s, ",")
	for _, p := range pairs {
		kv := strings.SplitN(p, ":", 2)
		if len(kv) == 2 {
			topics[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return topics
}
