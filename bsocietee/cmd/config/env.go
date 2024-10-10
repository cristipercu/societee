package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
  DbHost string
  DbUser string
  DbName string
  DbPassword string
  Sslmode string
  DbPort int64
  JWTExpirationInSeconds int64
  JWTSecret string
}

var Envs = initConfig()

func initConfig() Config {
  godotenv.Load()
  return Config{
    DbHost: getEnv("POSTGRESS_HOST", "localhost"),
    DbUser: getEnv("POSTGRES_USER", "postgres"),
    DbName: getEnv("PPOSTGRES_DB", "postgres"),
    DbPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
    Sslmode: getEnv("SSL_MODE", "disable"),
    DbPort: getEnvAsInt("POSTGRES_PORT", 5432),
    JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 60*60*24*7),
    // TODO: Do not fortget to add JWT_SECRET in your env, otherwise this default value is ised in you are cooked
    JWTSecret: getEnv("JWT_SECRET", "please_change_me!you_should_change_me!"),
  }
}

func getEnv(key, fallback string) string {
  if value, ok := os.LookupEnv(key); ok {
    return value
  }
  return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
  if value, ok := os.LookupEnv(key); ok {
    i, err := strconv.ParseInt(value, 10, 64)
    if err != nil {
      return fallback
    }
    return i
  }
  return fallback
}
