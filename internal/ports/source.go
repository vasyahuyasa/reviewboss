package ports

import "github.com/vasyahuyasa/reviewboss/internal/domain"

type MergeRequestSource interface {
	ListOpen() ([]*domain.MergeRequest, error)
	Get(id string) (*domain.MergeRequest, error)
	Save(mr *domain.MergeRequest) error
}
