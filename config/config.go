package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AWS
	Database
	API
	Kafka
	Redis
	Logger
}

type AWS struct {
	Region   string
	SqsTopic SQSTopic
}

type SQSTopic struct {
	Account string
}

type Database struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	Driver          string
}

type API struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Kafka struct {
	Brokers []string
	// ProducerTopic is where this service publishes Account events for other
	// services to consume.
	ProducerTopic string
	// ConsumerTopic is another service's event stream. Keep it different from
	// ProducerTopic or this consumer group would replay its own output.
	ConsumerTopic string
	ConsumerGroup string
}

type Redis struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func (r Redis) Addr() string {
	return r.Host + ":" + r.Port
}

func LoadConfig() *Config {
	return &Config{
		AWS: AWS{
			Region: GetEnv("AWS_REGION", "ap-southeast-1"),
			SqsTopic: SQSTopic{
				Account: GetEnv("ACCOUNTING_SQS", ""),
			},
		},
		Database: Database{
			Host:            GetEnv("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 3306),
			Username:        GetEnv("DB_USERNAME", "app"),
			Password:        GetEnv("DB_PASSWORD", ""),
			Database:        GetEnv("DB_NAME", "application"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME_SEC", 1800)) * time.Second,
			Driver:          GetEnv("DB_DRIVER", "mysql"),
		},
		API: API{
			Port:         GetEnv("API_PORT", "8080"),
			ReadTimeout:  time.Duration(getEnvInt("API_READ_TIMEOUT_SEC", 30)) * time.Second,
			WriteTimeout: time.Duration(getEnvInt("API_WRITE_TIMEOUT_SEC", 30)) * time.Second,
		},
		Kafka: Kafka{
			Brokers:       getEnvStringSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			ProducerTopic: GetEnv("KAFKA_PRODUCER_TOPIC", "account.events"),
			ConsumerTopic: GetEnv("KAFKA_CONSUMER_TOPIC", "identity.user-events"),
			ConsumerGroup: GetEnv("KAFKA_CONSUMER_GROUP", "account-service"),
		},
		Redis: Redis{
			Host:     GetEnv("REDIS_HOST", "localhost"),
			Port:     GetEnv("REDIS_PORT", "6379"),
			Password: GetEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Logger: loadLoggerConfig(),
	}
}

func GetEnv(key string, defaultValue string) string {
	result := os.Getenv(key)
	if result != "" {
		return result
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return n
}

func getEnvBool(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}
	return b
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultValue
	}
	return d
}

func getEnvStringSlice(key string, defaultValue []string) []string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	parts := strings.Split(val, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
