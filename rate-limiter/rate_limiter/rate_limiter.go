package limiter

type RateLimiter interface {
	IsRequestAllowed(clientId string) bool
}
