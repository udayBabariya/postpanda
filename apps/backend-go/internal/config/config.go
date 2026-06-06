package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	FrontendURL string
	BackendURL  string

	// Auth providers
	GitHubClientID     string
	GitHubClientSecret string
	GoogleClientID     string
	GoogleClientSecret string

	// Social media platform credentials
	LinkedInClientID     string
	LinkedInClientSecret string
	FacebookClientID     string
	FacebookClientSecret string
	InstagramClientID    string
	TikTokClientID       string
	TikTokClientSecret   string
	RedditClientID       string
	RedditClientSecret   string
	PinterestClientID    string
	PinterestClientSecret string
	YoutubeClientID      string
	YoutubeClientSecret  string
	SlackClientID        string
	SlackClientSecret    string
	DiscordClientID      string
	DiscordClientSecret  string
	TelegramBotToken     string

	// Storage
	CloudflareAccountID       string
	CloudflareAccessKey       string
	CloudflareSecretAccessKey string
	CloudflareBucketName      string
	StorageProvider           string
	UploadDirectory           string

	// Stripe
	StripeSecretKey      string
	StripeWebhookSecret  string
	StripeMonthlyPriceID string
	StripeYearlyPriceID  string

	// Email (Resend)
	ResendAPIKey string
	EmailFrom    string

	// OpenAI
	OpenAIAPIKey string

	// Misc
	IsGenericOAuthEnabled bool
	GenericOAuthURL       string
	GenericOAuthClientID  string
	GenericOAuthSecret    string

	// Rate limiting
	ThrottleLimit int
	ThrottleTTL   int

	// Upload
	MaxUploadSize int64

	// Sentry
	SentryDSN string

	// Encryption
	EncryptionKey string

	// Self-hosted
	DisableRegistration bool
	AllowedEmails       string

	// NX Cloud / temporal
	TemporalAddress string
	TemporalTaskQueue string
}

var App *Config

func Load() *Config {
	App = &Config{
		Port:        getEnv("PORT", "3000"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "secret"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:4200"),
		BackendURL:  getEnv("BACKEND_URL", "http://localhost:3000"),

		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),

		LinkedInClientID:      getEnv("LINKEDIN_CLIENT_ID", ""),
		LinkedInClientSecret:  getEnv("LINKEDIN_CLIENT_SECRET", ""),
		FacebookClientID:      getEnv("FACEBOOK_CLIENT_ID", ""),
		FacebookClientSecret:  getEnv("FACEBOOK_CLIENT_SECRET", ""),
		TikTokClientID:        getEnv("TIKTOK_CLIENT_ID", ""),
		TikTokClientSecret:    getEnv("TIKTOK_CLIENT_SECRET", ""),
		RedditClientID:        getEnv("REDDIT_CLIENT_ID", ""),
		RedditClientSecret:    getEnv("REDDIT_CLIENT_SECRET", ""),
		PinterestClientID:     getEnv("PINTEREST_CLIENT_ID", ""),
		PinterestClientSecret: getEnv("PINTEREST_CLIENT_SECRET", ""),
		YoutubeClientID:       getEnv("YOUTUBE_CLIENT_ID", ""),
		YoutubeClientSecret:   getEnv("YOUTUBE_CLIENT_SECRET", ""),
		SlackClientID:         getEnv("SLACK_CLIENT_ID", ""),
		SlackClientSecret:     getEnv("SLACK_CLIENT_SECRET", ""),
		DiscordClientID:       getEnv("DISCORD_CLIENT_ID", ""),
		DiscordClientSecret:   getEnv("DISCORD_CLIENT_SECRET", ""),
		TelegramBotToken:      getEnv("TELEGRAM_BOT_TOKEN", ""),

		CloudflareAccountID:       getEnv("CLOUDFLARE_ACCOUNT_ID", ""),
		CloudflareAccessKey:       getEnv("CLOUDFLARE_ACCESS_KEY", ""),
		CloudflareSecretAccessKey: getEnv("CLOUDFLARE_SECRET_ACCESS_KEY", ""),
		CloudflareBucketName:      getEnv("CLOUDFLARE_BUCKET_NAME", ""),
		StorageProvider:           getEnv("STORAGE_PROVIDER", "local"),
		UploadDirectory:           getEnv("UPLOAD_DIRECTORY", "./uploads"),

		StripeSecretKey:     getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),

		ResendAPIKey: getEnv("RESEND_API_KEY", ""),
		EmailFrom:    getEnv("EMAIL_FROM", "noreply@postpanda.io"),

		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),

		MaxUploadSize: int64(getEnvInt("MAX_UPLOAD_SIZE", 52428800)), // 50MB

		SentryDSN:     getEnv("SENTRY_DSN", ""),
		EncryptionKey: getEnv("ENCRYPTION_KEY", ""),

		DisableRegistration: getEnvBool("DISABLE_REGISTRATION", false),
		AllowedEmails:       getEnv("ALLOWED_EMAILS", ""),

		TemporalAddress:   getEnv("TEMPORAL_ADDRESS", "localhost:7233"),
		TemporalTaskQueue: getEnv("TEMPORAL_TASK_QUEUE", "main"),

		ThrottleLimit: getEnvInt("THROTTLE_LIMIT", 90),
		ThrottleTTL:   getEnvInt("THROTTLE_TTL", 3600),

		IsGenericOAuthEnabled: getEnvBool("IS_GENERIC_OAUTH_ENABLED", false),
		GenericOAuthURL:       getEnv("GENERIC_OAUTH_URL", ""),
		GenericOAuthClientID:  getEnv("GENERIC_OAUTH_CLIENT_ID", ""),
		GenericOAuthSecret:    getEnv("GENERIC_OAUTH_SECRET", ""),
	}
	return App
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return defaultVal
}
