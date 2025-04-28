package outbound

import (
	"github.com/vasyahuyasa/reviewboss/internal/domain"
	"github.com/vasyahuyasa/reviewboss/internal/ports"
)

type GitLabSource struct { /* client, config */
}

func (g *GitLabSource) ListOpen() ([]*domain.MergeRequest, error)   { /* call GitLab API */ }
func (g *GitLabSource) Get(id string) (*domain.MergeRequest, error) { /* ... */ }
func (g *GitLabSource) Save(mr *domain.MergeRequest) error          { /* ... */ }

var _ ports.MergeRequestSource = (*GitLabSource)(nil)
