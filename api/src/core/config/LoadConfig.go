package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppName:              os.Getenv("APP_NAME"),
		AppPort:              os.Getenv("APP_PORT"),
		DBHost:               os.Getenv("DB_HOST"),
		DBPort:               os.Getenv("DB_PORT"),
		DBUser:               os.Getenv("DB_USER"),
		DBPassword:           os.Getenv("DB_PASSWORD"),
		DBName:               os.Getenv("DB_NAME"),
		DBSSL:                os.Getenv("DB_SSL"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		GoogleOAuthClientIDs: os.Getenv("GOOGLE_OAUTH_CLIENT_IDS"),
		CORSOrigin:           os.Getenv("CORS_ORIGIN"),
		CookieSecure:         os.Getenv("COOKIE_SECURE") == "true",
		CloudinaryCloudName:  os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:     os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret:  os.Getenv("CLOUDINARY_API_SECRET"),
		RabbitMQURL:          os.Getenv("RABBITMQ_URL"),
		NLPServiceURL:        os.Getenv("NLP_SERVICE_URL"),
		// Contenido JSON completo de la cuenta de servicio (Firebase Console
		// -> Configuración -> Cuentas de servicio -> Generar nueva clave
		// privada), no una ruta de archivo -- más simple de configurar en
		// Railway/Heroku que subir un archivo al contenedor. Opcional: si
		// falta, el envío de push queda en no-op (ver src/core/push).
		FirebaseServiceAccountKey: os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY"),
	}

	if cfg.AppPort == "" {
		cfg.AppPort = os.Getenv("PORT")
	}
	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}
	if cfg.CORSOrigin == "" {
		cfg.CORSOrigin = "*"
		// Con "*" el middleware CORS (ver core/middleware/CORS.go) no manda
		// Allow-Credentials -- cualquier login/flujo que dependa de la
		// cookie de sesión desde un navegador no va a funcionar hasta que
		// se configure el origen real acá.
		log.Println("advertencia: CORS_ORIGIN no configurado, usando \"*\" -- el login por cookie desde un navegador no funcionará hasta que se fije un origen real")
	}
	if cfg.RabbitMQURL == "" {
		cfg.RabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}
	if cfg.NLPServiceURL == "" {
		cfg.NLPServiceURL = "http://localhost:8006"
	}
	if cfg.GoogleOAuthClientIDs == "" {
		// Client ID "web" (client_type 3) del proyecto de Firebase en
		// android/app/google-services.json del app Flutter -- no es un
		// secreto (viaja embebido en el apk), así que sirve de default
		// razonable si no se sobreescribe con la variable de entorno.
		cfg.GoogleOAuthClientIDs = "49495820436-oq98o62p3vjedpee1v1k6afqadnt2i2o.apps.googleusercontent.com"
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
