package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppName:             os.Getenv("APP_NAME"),
		AppPort:             os.Getenv("APP_PORT"),
		DBHost:              os.Getenv("DB_HOST"),
		DBPort:              os.Getenv("DB_PORT"),
		DBUser:              os.Getenv("DB_USER"),
		DBPassword:          os.Getenv("DB_PASSWORD"),
		DBName:              os.Getenv("DB_NAME"),
		DBSSL:               os.Getenv("DB_SSL"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		CORSOrigin:          os.Getenv("CORS_ORIGIN"),
		CookieSecure:        os.Getenv("COOKIE_SECURE") == "true",
		CloudinaryCloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret: os.Getenv("CLOUDINARY_API_SECRET"),
	}

	if cfg.AppPort == "" {
		cfg.AppPort = os.Getenv("PORT")
	}
	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}
	if cfg.CORSOrigin == "" {
		cfg.CORSOrigin = "*"
	}

	required := map[string]string{
		"DB_HOST":               cfg.DBHost,
		"DB_PORT":               cfg.DBPort,
		"DB_USER":               cfg.DBUser,
		"DB_PASSWORD":           cfg.DBPassword,
		"DB_NAME":               cfg.DBName,
		"JWT_SECRET":            cfg.JWTSecret,
		"CLOUDINARY_CLOUD_NAME": cfg.CloudinaryCloudName,
		"CLOUDINARY_API_KEY":    cfg.CloudinaryAPIKey,
		"CLOUDINARY_API_SECRET": cfg.CloudinaryAPISecret,
	}
	for name, value := range required {
		if value == "" {
			return nil, fmt.Errorf("falta la variable de entorno obligatoria: %s", name)
		}
	}

	return cfg, nil
}
