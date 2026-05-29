package mixer

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	StemsDir          string
	OutputDir         string
	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2StemsBucket     string
)

func init() {
	// Best-effort: load .env from working directory (same as Python's load_dotenv())
	_ = godotenv.Load()

	StemsDir = getEnv("STEMS_DIR", "/Users/rca/nltl/lvg-bucket/mp3/first-principles")
	OutputDir = getEnv("OUTPUT_DIR", "./output")
	R2AccountID = os.Getenv("R2_ACCOUNT_ID")
	R2AccessKeyID = os.Getenv("R2_ACCESS_KEY_ID")
	R2SecretAccessKey = os.Getenv("R2_SECRET_ACCESS_KEY")
	R2StemsBucket = os.Getenv("R2_STEMS_BUCKET")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
