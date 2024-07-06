package circuit_breaker

import (
	"errors"
	"time"
)

type Request struct {
	exp func() (*Response, error)
}

type Response struct {
	statusCode int
}

type CircuitBreaker struct {
	maxFailuresAllowed int
	failed             int
	circuitOpen        bool
	circuitOpenedAt    time.Time
}

func NewCircuitBreaker(mxFail int) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailuresAllowed: mxFail,
	}
}

func (s *CircuitBreaker) Check(request Request) (*Response, error) {
	if s.circuitOpen {
		// check if circuit is opened for more than 10 minutes
		if time.Now().After(s.circuitOpenedAt.Add(10 * time.Minute)) {
			s.failed = 0
			s.circuitOpen = false
		} else {
			return nil, errors.New("circuit is opened")
		}
	}
	result, err := request.exp()
	if err != nil {
		s.failed++
		if s.failed > s.maxFailuresAllowed {
			s.circuitOpen = true
			s.circuitOpenedAt = time.Now()
		}
		return nil, err
	}
	return result, nil
}
