package main

type Notificator interface {
	NotifyReviewerProposal(mr *MergeRequest, reviewer *Reviewer) error
	NotifyNoReviewersAvailable(mr *MergeRequest) error
	NotifyMRLongInReview(mr *MergeRequest) error
	NotifyMRReadyToMerge(mr *MergeRequest) error
}
