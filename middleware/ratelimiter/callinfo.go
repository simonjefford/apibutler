package ratelimiter

import "time"

type CallInfo struct {
	Count   int
	Start   time.Time
	Limit   int
	Seconds time.Duration
}

func (c *CallInfo) Remaining() int {
	return c.Limit - c.Count
}

func (c *CallInfo) timeSinceLastReset() time.Duration {
	return time.Since(c.Start)
}

func (c *CallInfo) IsLimitExceeded() bool {
	dur := c.timeSinceLastReset()
	return dur < c.Seconds && c.Count >= c.Limit
}

func (c *CallInfo) ResetIfNeccesary() {
	dur := c.timeSinceLastReset()
	if dur > c.Seconds {
		c.Start = time.Now()
		c.Count = 0
	}
}

func NewCallInfo(limit int, seconds int) *CallInfo {
	return &CallInfo{
		Start:   time.Now(),
		Limit:   limit,
		Seconds: time.Duration(seconds) * time.Second,
	}
}
