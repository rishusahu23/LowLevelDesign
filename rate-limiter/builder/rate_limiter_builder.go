package builder

import limiter "github.com/rishu/design/rate-limiter/rate_limiter"

type RateLimiterBuilder interface {
	SetRate(rate int) RateLimiterBuilder
	SetBucketSize(size int) RateLimiterBuilder
	Build() limiter.RateLimiter
}
