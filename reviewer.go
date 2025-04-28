package main

type ReviewerService interface {
	ProposeReviewer(mr *MergeRequest) (*Reviewer, error)
	AssignReviewer(mr *MergeRequest, reviewer *Reviewer) error
	DeclineReview(mr *MergeRequest, reviewer *Reviewer) error
	NoReviewersAvailable(mr *MergeRequest) error
}
