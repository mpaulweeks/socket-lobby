package main

import (
	"time"
)

type Clock interface {
	Now() time.Time
	NowTicks() int
}

type realClock struct{}

func (c *realClock) Now() time.Time {
	return time.Now()
}

func (c *realClock) NowTicks() int {
	return int(c.Now().Unix())
}

var RealClock Clock = new(realClock)
