package configs

import (
	"os"
	"strings"

	"github.com/ShopOnGO/ShopOnGO/configs"
	"github.com/joho/godotenv"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
)

type Config struct {
	Db            DbConfig
	KafkaProduct  KafkaConsumerConfig
	KafkaVariant  KafkaConsumerConfig
	KafkaMedia    KafkaConsumerConfig
	KafkaProducer KafkaProducerConfig
	LogLevel      logger.LogLevel
	FileLogLevel  logger.LogLevel
}

type DbConfig struct {
	Dsn string
}

type KafkaConsumerConfig struct {
	Brokers  []string
	Topic    string
	GroupID  string
	ClientID string
}

type KafkaProducerConfig struct {
	Brokers []string
	Topic   map[string]string
}

func LoadConfig() *Config {
	if _, err := os.Stat(".env"); err == nil {
		// Локально есть .env → загружаем
		if loadErr := godotenv.Load(); loadErr != nil {
			logger.Error("Failed to load .env file", loadErr.Error())
		}
	} else {
		// В контейнере файла нет → просто идём дальше
		logger.Info(".env not found, using environment variables only")
	}

	brokersRaw := os.Getenv("KAFKA_BROKERS")
	brokers := strings.Split(brokersRaw, ",")
	// logger
	logLevelStr := os.Getenv("PRODUCT_SERVICE_LOG_LEVEL")
	if logLevelStr == "" {
		logLevelStr = "INFO"
	}
	LogLevel := configs.ParseLogLevel(logLevelStr)
	fileLogLevelStr := os.Getenv("PRODUCT_SERVICE_FILE_LOG_LEVEL")
	if fileLogLevelStr == "" {
		fileLogLevelStr = "INFO"
	}
	FileLogLevel := configs.ParseLogLevel(fileLogLevelStr)

	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		KafkaProduct: KafkaConsumerConfig{
			Brokers:  brokers,
			Topic:    os.Getenv("KAFKA_PRODUCT_TOPIC"),
			GroupID:  os.Getenv("KAFKA_PRODUCT_GROUP_ID"),
			ClientID: os.Getenv("KAFKA_PRODUCT_CLIENT_ID"),
		},
		KafkaVariant: KafkaConsumerConfig{
			Brokers:  brokers,
			Topic:    os.Getenv("KAFKA_VARIANT_TOPIC"),
			GroupID:  os.Getenv("KAFKA_VARIANT_GROUP_ID"),
			ClientID: os.Getenv("KAFKA_VARIANT_CLIENT_ID"),
		},
		KafkaMedia: KafkaConsumerConfig{
			Brokers:  brokers,
			Topic:    os.Getenv("KAFKA_MEDIA_TOPIC"),
			GroupID:  os.Getenv("KAFKA_MEDIA_GROUP_ID"),
			ClientID: os.Getenv("KAFKA_MEDIA_CLIENT_ID"),
		},
		KafkaProducer: KafkaProducerConfig{
			Brokers: brokers,
			Topic:   parseKafkaTopics(os.Getenv("KAFKA_PRODUCER_TOPIC")),
		},
		LogLevel:     LogLevel,
		FileLogLevel: FileLogLevel,
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
