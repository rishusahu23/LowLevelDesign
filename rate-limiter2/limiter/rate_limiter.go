package limiter

type RateLimiter interface {
	IsRequestAllowed(apiKey string) bool
}
