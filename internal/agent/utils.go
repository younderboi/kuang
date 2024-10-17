package agent

import (
	"math/rand"
	"time"
)

func (agent *Agent) backoff(retries int) time.Duration {
	// Compute the backoff delay: baseDelay * 2^retries + jitter
	// The jitter is a random duration added to prevent thundering herd problems
	jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
	return agent.BaseDelay*(1<<retries) + jitter
}
