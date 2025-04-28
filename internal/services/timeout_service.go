package services

import (
	"time"

	"github.com/vasyahuyasa/reviewboss/internal/domain"
	"github.com/vasyahuyasa/reviewboss/internal/ports"
)

type TimeoutService struct {
	Checker ports.TimeoutChecker
}

func (ts *TimeoutService) RunPeriodic(mrs []*domain.MergeRequest) {
	now := time.Now().Unix()
	for _, mr := range mrs {
		ts.Checker.Check(mr)
	}
}
