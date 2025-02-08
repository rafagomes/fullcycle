package domain

import "time"

type Result struct {
	StatusCode int
	Err        error
}

type Report struct {
	TotalTime          time.Duration
	TotalRequests      int
	Status200Count     int
	ErrorCount         int
	StatusDistribution map[int]int
}
