package config

import "os"

type Config struct {
	ServerAddress string
	DBUsername    string
	DBPassword    string
	DBHost        string
	DBName        string
	DBPort        string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string
	AccessKey     string
	JWTSecretKey  string
	JWTRefreshKey string
	SMTPServer    string
	SMTPPort      string
	SenderEmail   string
	EmailPassword string
}

func LoadConfig() *Config {
	return &Config{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		DBUsername:    os.Getenv("MYSQL_USER"),
		DBPassword:    os.Getenv("MYSQL_PASSWORD"),
		DBHost:        os.Getenv("DB_HOST"),
		DBName:        os.Getenv("MYSQL_DATABASE"),
		DBPort:        os.Getenv("DB_PORT"),

		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       os.Getenv("REDIS_DB"),
		AccessKey:     os.Getenv("ACCESS_KEY"),
		JWTSecretKey:  os.Getenv("JWT_SECRET_KEY"),
		JWTRefreshKey: os.Getenv("JWT_REFRESH_KEY"),
		SMTPServer:    os.Getenv("SMTP_SERVER"),
		SMTPPort:      os.Getenv("SMTP_PORT"),
		SenderEmail:   os.Getenv("SEND_EMAIL"),
		EmailPassword: os.Getenv("EMAIL_PASSWORD"),
	}
}
