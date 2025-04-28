package ports

import "github.com/vasyahuyasa/reviewboss/internal/domain"

type Notifier interface {
	NotifyProposal(mr *domain.MergeRequest, reviewer string) error
	NotifyNoReviewer(mr *domain.MergeRequest) error
	NotifyLongInReview(mr *domain.MergeRequest) error
	NotifyMerged(mr *domain.MergeRequest) error
}
