package builder

import limiter "github.com/rishu/design/rate-limiter/rate_limiter"

type TokenBucketBuilder struct {
	Rate       int
	BucketSize int
}

func NewTokenBucketBuilder() *TokenBucketBuilder {
	return &TokenBucketBuilder{}
}

func (b *TokenBucketBuilder) SetRate(rate int) RateLimiterBuilder {
	b.Rate = rate
	return b
}

func (b *TokenBucketBuilder) SetBucketSize(size int) RateLimiterBuilder {
	b.BucketSize = size
	return b
}

func (b *TokenBucketBuilder) Build() limiter.RateLimiter {
	return limiter.NewTokenBucket(b.Rate, b.BucketSize)
}
