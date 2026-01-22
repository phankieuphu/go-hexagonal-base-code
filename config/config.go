package config

import (
	"os"
	"time"
)

type Config struct {
	AWS
	Database
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

func LoadConfig() *Config {

	return &Config{
		AWS: AWS{
			Region: GetEnv("AWS_REGION", "ap-southeast-1"),
			SqsTopic: SQSTopic{
				Account: GetEnv("ACCOUNTING_SQS", ""),
			},
		},
		Database: Database{
			// AccountTableName: GetEnv("ACCOUNT_ENTRY_TABLE", "accounting_account_entry"),
		},
	}
}

func GetEnv(key string, defaultValue string) string {
	result := os.Getenv(key)
	if result != "" {
		return result
	}
	return defaultValue

}
