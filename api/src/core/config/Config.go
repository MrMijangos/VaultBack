package config

type Config struct {
	AppName             string
	AppPort             string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBSSL               string
	JWTSecret           string
	CORSOrigin          string
	CookieSecure        bool
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
	RabbitMQURL         string
	NLPServiceURL       string
}
