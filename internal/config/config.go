package config

import (
	"errors"
	"fmt"
	"github.com/caarlos0/env/v9"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"net"
	"strings"
)

type Config struct {
	Port      int    `env:"PORT,required" validate:"required,min=1,max=65535"`
	Host      string `env:"HOST,required" validate:"required,hostname_or_ip"`
	JwtSecret string `env:"JWT_SECRET,required" validate:"required"`
	DbHost    string `env:"DB_HOST,required" validate:"required,hostname_or_ip"`
	DbPort    int    `env:"DB_PORT,required" validate:"required,min=1,max=65535"`
	DbUser    string `env:"DB_USER"`
	DbPass    string `env:"DB_PASS"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Файл .env не найден, загрузка переменных из окружения")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateConfig(cfg *Config) error {
	validate := validator.New()

	_ = validate.RegisterValidation("hostname_or_ip", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		return net.ParseIP(value) != nil || strings.Contains(value, ".")
	})

	err := validate.Struct(cfg)
	if err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("Ошибка в поле %s: %s", err.Field(), err.Tag()))
		}
		return errors.New(strings.Join(validationErrors, ", "))
	}

	return nil
}
