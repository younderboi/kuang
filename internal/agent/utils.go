package agent

import (
	"math/rand"
	"time"
)

func (agent *Agent) backoff(retries int) time.Duration {
	// Compute re-connect interval using exponential backoff:
	// (baseDelay * 2^retries) + jitter
	jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
	return agent.baseDelay*(1<<retries) + jitter
}
