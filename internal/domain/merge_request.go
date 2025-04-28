package domain

import "time"

type MRState string

const (
	StateNew                   MRState = "new_mr"
	StateWaitingVoluntary      MRState = "waiting_voluntary"
	StateVoluntaryExpired      MRState = "voluntary_expired"
	StateProposeReviewer       MRState = "propose_reviewer"
	StateReviewerDeclined      MRState = "reviewer_declined"
	StateNoReviewers           MRState = "no_reviewers"
	StateWaitingSomeone        MRState = "waiting_someone"
	StateWaitingSomeoneExpired MRState = "waiting_someone_expired"
	StateAgreementReceived     MRState = "agreement_received"
	StateProposalExpired       MRState = "proposal_expired"
	StateAssignReviewer        MRState = "assign_reviewer"
	StateReviewInProgress      MRState = "review_in_progress"
	StateLongInReview          MRState = "long_in_review"
	StateMerged                MRState = "merged"
)

type MergeRequest struct {
	ID               string
	RepoID           string
	Author           string
	State            MRState
	AssignedReviewer string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
