package ports

import "github.com/vasyahuyasa/reviewboss/internal/domain"

type TimeoutChecker interface {
	Check(mr *domain.MergeRequest) error
}
