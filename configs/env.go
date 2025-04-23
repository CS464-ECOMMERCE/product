package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v81"
)

var (
	PORT            string
	API_LISTEN_HOST string

	POSTGRESQL_CONN_STRING_MASTER string
	POSTGRESQL_CONN_STRING_SLAVE  string
	POSTGRESQL_MAX_IDLE_CONNS     int
	POSTGRESQL_MAX_OPEN_CONNS     int
	PRODUCT_SERVICE_ADDR          string
	CART_SERVICE_ADDR             string
	S3_ENDPOINT                   string
	S3_ACCESS_KEY                 string
	S3_SECRET_KEY                 string
	S3_BUCKET                     string
	GCS_CREDENTIALS               string
	GCS_PROJECT_ID                string
	FRONTEND_URL                  string
	GOOGLE_MAPS_API_KEY           string
	ENVIRONMENT                   string
)

func InitEnv() {
	// loads environment variables
	envPath := "/app/secrets/.env"
	if os.Getenv("ENV") == "dev" {
		envPath = "./secrets/testing.env"
	}
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Println(err)
		panic("Error loading env file")
	}

	// rest api
	PORT = getEnv("API_PORT", "8080")
	API_LISTEN_HOST = getEnv("API_LISTEN_HOST", "0.0.0.0")

	// postgress
	POSTGRESQL_CONN_STRING_MASTER = getEnv("POSTGRESQL_CONN_STRING_MASTER", "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai")
	POSTGRESQL_CONN_STRING_SLAVE = getEnv("POSTGRESQL_CONN_STRING_SLAVE", "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai")
	maxOpenConns, err := strconv.Atoi(getEnv("POSTGRESQL_MAX_OPEN_CONNS", "10"))
	if err != nil {
		panic("Invalid value for POSTGRESQL_MAX_OPEN_CONNS")
	}
	POSTGRESQL_MAX_OPEN_CONNS = maxOpenConns
	maxIdleConns, err := strconv.Atoi(getEnv("POSTGRESQL_MAX_IDLE_CONNS", "5"))
	if err != nil {
		panic("Invalid value for POSTGRESQL_MAX_IDLE_CONNS")
	}
	POSTGRESQL_MAX_IDLE_CONNS = maxIdleConns

	// grpc
	PRODUCT_SERVICE_ADDR = getEnv("PRODUCT_SERVICE_ADDR", "product.default.svc.cluster.local:50051")
	CART_SERVICE_ADDR = getEnv("CART_SERVICE_ADDR", "cart.default.svc.cluster.local:50051")

	//S3
	S3_ENDPOINT = getEnv("S3_ENDPOINT", "minio.default.svc.cluster.local:9000")
	S3_ACCESS_KEY = getEnv("S3_ACCESS_KEY", "minioadmin")
	S3_SECRET_KEY = getEnv("S3_SECRET_KEY", "minioadmin")
	S3_BUCKET = getEnv("S3_BUCKET", "cs464-application")
	GCS_CREDENTIALS = "/app/secrets/backend-sa-key.json"
	GCS_PROJECT_ID = getEnv("GCS_PROJECT_ID", "cs464-ecommerce")

	// Init Stripe
	stripe.Key = getEnv("STRIPE_SECRET_KEY", "some-secret-key")

	// frontend (for stripe)
	FRONTEND_URL = getEnv("FRONTEND_URL", "http://localhost:3000")

	// Google Maps API Key
	GOOGLE_MAPS_API_KEY = getEnv("GOOGLE_MAPS_API_KEY", "your-google-maps-api-key")

	// Environment
	ENVIRONMENT = getEnv("ENVIRONMENT", "dev")
}

// get env with default if the value is empty
// getEnv("ENV_VAR", "default")
func getEnv(s ...string) string {

	if len(s) <= 0 {

		// only one arg, don't provide defaults
		return ""

	} else if val := os.Getenv(s[0]); len(s) >= 2 && val != "" {

		// two args and the env var provides empty value
		return val

	} else {

		return s[1]

	}

}
