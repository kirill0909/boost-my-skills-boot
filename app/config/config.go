package config

import (
	"aifory-pay-admin-bot/internal/utils"
	"fmt"
	"time"

	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Server            Server
	Postgres          Postgres
	Logger            Logger
	OpenTelemetry     OpenTelemetry
	AiforyPayAdminLog AiforyPayAdminLog
}

var C *Config

type Server struct {
	AppVersion      string        `validate:"required"`
	Host            string        `validate:"required"`
	GRPCPort        string        `validate:"required"`
	HTTPPort        string        `validate:"required"`
	TGToken         string        `validate:"required"`
	ShutdownTimeout time.Duration `validate:"required"`
}

type Postgres struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	DBName   string `validate:"required"`
	SSLMode  string `validate:"required"`
	PGDriver string `validate:"required"`
	Settings struct {
		MaxOpenConns    int           `validate:"required,min=1"`
		ConnMaxLifetime time.Duration `validate:"required,min=1"`
		MaxIdleConns    int           `validate:"required,min=1"`
		ConnMaxIdleTime time.Duration `validate:"required,min=1"`
	}
}

type Logger struct {
	Level          string `validate:"required"`
	SkipFrameCount int
	InFile         bool
	FilePath       string
	InTG           bool
	TGLevel        string `validate:"required"`
	ChatID         int64
	TGToken        string
	AlertUsers     []string
}

type OpenTelemetry struct {
	Host        string `validate:"required"`
	ServiceName string `validate:"required"`
}

type AiforyPayAdminLog struct {
	ChatID int64
	TagMe  []string
}

func loadConfig() (*viper.Viper, error) {
	v := viper.New()

	v.AddConfigPath(fmt.Sprintf("./%s", utils.ConfigPath))
	v.SetConfigName(utils.ConfigFileName)
	v.SetConfigType(utils.ConfigExtension)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, nil
}

func parseConfig(v *viper.Viper) (*Config, error) {
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

func ProvideConfig() (*Config, error) {
	cfgFile, err := loadConfig()
	if err != nil {
		log.Printf("LoadConfig: %s", err.Error())
		return nil, err
	}
	log.Println("Config loaded")

	cfg, err := parseConfig(cfgFile)
	if err != nil {
		log.Printf("ParseConfig: %s", err.Error())
		return nil, err
	}
	log.Print("Config parsed")
	C = cfg
	return cfg, nil
}
