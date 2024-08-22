package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kenmobility/github-api-hex/common/helpers"
	"gopkg.in/go-playground/validator.v9"
)

type Config struct {
	AppEnv           string
	GitHubToken      string `validate:"required"`
	DatabaseHost     string `validate:"required"`
	DatabasePort     string `validate:"required"`
	DatabaseUser     string `validate:"required"`
	DatabasePassword string `validate:"required"`
	DatabaseName     string `validate:"required"`
	FetchInterval    time.Duration
	GitHubApiBaseURL string
	DefaultStartDate time.Time
	DefaultEndDate   time.Time
	Address          string
	Port             string
}

func LoadConfig(path string) *Config {
	var err error

	if path == "" {
		path = ".env"
	}
	if err := godotenv.Load(path); err != nil {
		log.Fatal("env config error: ", err)
	}

	interval := os.Getenv("FETCH_INTERVAL")
	if interval == "" {
		interval = "1h"
	}

	intervalDuration, err := time.ParseDuration(interval)
	if err != nil {
		log.Fatalf("Invalid FETCH_INTERVAL :[%s] env format: %v", interval, err)
	}

	var sDate time.Time
	var eDate time.Time

	startDate := os.Getenv("DEFAULT_START_DATE")
	if startDate == "" {
		sDate = time.Now().AddDate(0, -4, 0)
	} else {
		sDate, err = time.Parse(time.RFC3339, startDate)
		if err != nil {
			log.Fatalf("Invalid DEFAULT_START_DATE [%s] env format: %v", startDate, err)
		}
	}

	endDate := os.Getenv("DEFAULT_START_DATE")
	if endDate == "" {
		eDate = time.Now()
	} else {
		eDate, err = time.Parse(time.RFC3339, endDate)
		if err != nil {
			log.Fatalf("Invalid DEFAULT_END_DATE [%s] env format: %v", endDate, err)
		}
	}

	configVar := Config{
		AppEnv:           helpers.Getenv("APP_ENV", "local"),
		GitHubToken:      os.Getenv("GIT_HUB_TOKEN"),
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
		DatabasePort:     os.Getenv("DATABASE_PORT"),
		DatabaseUser:     os.Getenv("DATABASE_USER"),
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		FetchInterval:    intervalDuration,
		DefaultStartDate: sDate,
		DefaultEndDate:   eDate,
		GitHubApiBaseURL: os.Getenv("GITHUB_API_BASE_URL"),
		Address:          helpers.Getenv("ADDRESS", "127.0.0.1"),
		Port:             helpers.Getenv("PORT", "5000"),
	}

	validate := validator.New()
	err = validate.Struct(configVar)
	if err != nil {
		log.Fatalf("env validation error: %s", err.Error())
	}

	return &configVar
}
