package director

import (
	"github.com/rishu/design/rate-limiter/builder"
	limiter "github.com/rishu/design/rate-limiter/rate_limiter"
)

type RateLimiterDirector struct {
	Builder builder.RateLimiterBuilder
}

func NewRateLimiterDirector(builder builder.RateLimiterBuilder) *RateLimiterDirector {
	return &RateLimiterDirector{
		Builder: builder,
	}
}

func (d *RateLimiterDirector) SetBuilder(builder builder.RateLimiterBuilder) {
	d.Builder = builder
}

func (d *RateLimiterDirector) Construct() limiter.RateLimiter {
	return d.Builder.SetRate(5).SetBucketSize(10).Build()
}
