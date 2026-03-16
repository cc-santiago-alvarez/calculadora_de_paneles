package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port            int
	MongoDBURI      string
	Environment     string
	PVGISBaseURL    string
	NASAPowerURL    string
	APITimeout      int // milliseconds
	APIMaxRetries   int
}

func Load() *Config {
	return &Config{
		Port:          getEnvInt("PORT", 3001),
		MongoDBURI:    getEnv("MONGODB_URI", "mongodb://root:12345abc@10.2.20.113:27017/calculadora_paneles?authSource=admin&directConnection=true"),
		Environment:   getEnv("NODE_ENV", "development"),
		PVGISBaseURL:  "https://re.jrc.ec.europa.eu/api/v5_2",
		NASAPowerURL:  "https://power.larc.nasa.gov/api/temporal/monthly/point",
		APITimeout:    10000,
		APIMaxRetries: 3,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
