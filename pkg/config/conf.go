package config

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	DefaultApiPort         = "4000"
	DefaultEnvironment     = "dev"
	DefaultAppName         = "forge-api"
	DefaultKeycloakBackend = "keycloak-backend-variable-not-found"
	RedisDefaultHost       = "redis:6379"
	RedisDefaultPassword   = ""
	RedisDefaultDb         = "0"
	AwsAccessKeyId         = "a"
	AwsSecretAccessKey     = "b"
)

type RedisConf struct {
	Host     string
	Password string
	Db       string
}

type Conf struct {
	AppName                  string
	RedisDB                  RedisConf
	APIPort                  string
	Environment              string
	KeycloakHost             string
	KeycloakClient           string
	KeycloakHostWithoutRealm string
	KeycloakRealm            string
	ClientSecret             string
	PostgreUrl               string
	NatsUrl                  string
	AwsAccessKeyId           string
	AwsSecretAccessKey       string
	MicrosoftAppId           string
	MicrosoftAppSecret       string
	MicrosoftTenantId        string
	PostgresHost             string
	PostgresPort             string
	PostgresUser             string
	PostgresPassword         string
	PostgresDatabase         string
	PostgresSsl              string
}

func New() *Conf {
	_ = godotenv.Load()

	conf := Conf{
		AppName: getEnv("APP_NAME", DefaultAppName),
		RedisDB: RedisConf{
			Host:     getEnv("REDIS_HOST", RedisDefaultHost),
			Password: getEnv("REDIS_PASSWORD", RedisDefaultPassword),
			Db:       getEnv("REDIS_DB", RedisDefaultDb),
		},
		PostgreUrl:               getEnv("POSTGRES_URL", ""),
		APIPort:                  getEnv("API_PORT", DefaultApiPort),
		Environment:              getEnv("API_ENV", DefaultEnvironment),
		KeycloakHost:             getEnv("KEYCLOAK_API_HOST", DefaultKeycloakBackend),
		KeycloakHostWithoutRealm: getEnv("KEYCLOAK_HOST_WITHOUT_REALM", DefaultKeycloakBackend),
		KeycloakClient:           getEnv("KEYCLOAK_CLIENT", DefaultKeycloakBackend),
		KeycloakRealm:            getEnv("KEYCLOAK_REALM", ""),
		ClientSecret:             getEnv("CLIENT_SECRET", DefaultKeycloakBackend),
		NatsUrl:                  getEnv("NATS_URL", ""),
		AwsAccessKeyId:           getEnv("AWS_KEY", AwsAccessKeyId),
		AwsSecretAccessKey:       getEnv("AWS_SECRET_KEY", AwsSecretAccessKey),
		MicrosoftAppId:           getEnv("MICROSOFT_APP_ID", ""),
		MicrosoftAppSecret:       getEnv("MICROSOFT_CLIENT_SECRET", ""),
		MicrosoftTenantId:        getEnv("MICROSOFT_TENANT_ID", ""),
		PostgresHost:             getEnv("POSTGRES_HOST", ""),
		PostgresPort:             getEnv("POSTGRES_PORT", ""),
		PostgresUser:             getEnv("POSTGRES_USER", ""),
		PostgresPassword:         getEnv("POSTGRES_PASSWORD", ""),
		PostgresDatabase:         getEnv("POSTGRES_DATABASE", ""),
		PostgresSsl:              getEnv("POSTGRES_SSL", ""),
	}

	return &conf
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
