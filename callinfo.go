package main

import (
	"log"
	"time"
)

type CallInfo struct {
	Count int
	Start time.Time
}

func (c *CallInfo) timeSinceLastReset() time.Duration {
	return time.Since(c.Start)
}

func (c *CallInfo) LimitExceeded() bool {
	dur := c.timeSinceLastReset()
	return dur < time.Second*20 && c.Count >= 5
}

func (c *CallInfo) ResetIfNeccesary() {
	dur := c.timeSinceLastReset()
	log.Println(dur)
	if dur > time.Second*20 {
		log.Println("Resetting")
		c.Start = time.Now()
		c.Count = 0
	} else {
		log.Println("Not resetting")
	}
}

func NewCallInfo() *CallInfo {
	return &CallInfo{Start: time.Now()}
}
