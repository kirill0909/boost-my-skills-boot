package config

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Server                       Server
	Postgres                     Postgres
	Redis                        Redis
	TgBot                        TgBot
	AdminChatID                  int64  `validate:"required"`
	AwaitingDirectionNameDelay   int    `validate:"required"`
	AwaitingParentDirectionDelay int    `validate:"required"`
	AwaitingAddInfoDirection     int    `validate:"required"`
	GRPCApiKey                   string `validate:"required"`
}

type Server struct {
	HTTP struct {
		Host string `validate:"required"`
		Port string `validate:"required"`
	}
	GRPC struct {
		Host string `validate:"required"`
		Port string `validate:"required"`
	}
}

type Postgres struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	DBName   string `validate:"required"`
	SSLMode  string `validate:"required"`
}

type Redis struct {
	Host               string `validate:"required"`
	Port               string `validate:"required"`
	UserName           string `validate:"required"`
	Password           string `validate:"required"`
	MinIdleConns       int    `validate:"required"`
	PoolTimeout        int    `validate:"required"`
	DB                 int
	PoolSize           int
	UseCertificates    bool
	InsecureSkipVerify bool
	CertificatesPaths  struct {
		Cert string
		Key  string
		Ca   string
	}
}

type TgBot struct {
	ApiKey string `validate:"required"`
	Prefix string `validate:"required"`
}

func LoadConfig() (*viper.Viper, error) {
	v := viper.New()

	v.AddConfigPath(fmt.Sprintf("./%s", ConfigPath))
	v.SetConfigName(ConfigFileName)
	v.SetConfigType(ConfigExtension)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
		return nil, err
	}
	err = validator.New().Struct(c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
