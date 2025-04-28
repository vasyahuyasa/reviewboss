package services

import (
	"github.com/vasyahuyasa/reviewboss/internal/domain"
	"github.com/vasyahuyasa/reviewboss/internal/ports"
)

type ReviewerService struct {
	Store    ports.ReviewerStore
	Notifier ports.Notifier
}

func (s *ReviewerService) Propose(mr *domain.MergeRequest) error {
	cand, _ := s.Store.ListEligible(mr.RepoID, mr.Author)
	if len(cand) == 0 {
		s.Notifier.NotifyNoReviewer(mr)
		return nil
	}
	// pick first or random
	r := cand[0]
	s.Notifier.NotifyProposal(mr, r)
	return nil
}
