package domain

import (
	"time"
)

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

type Rewive struct{}

type Reviewer struct{}

type NotificationChannel interface {
	BroadcastMrRegistered() error
	BroadcastProposedReviewer(Reviewer) error
	BroadcastNoReviewers() error
	BroadcastAssignReviewer(Reviewer) error
	BroadcastReviewStatus() error
}

type MrRemoteStatus string

const (
	StatusOpened MrRemoteStatus = "opened"
	StatusMerged MrRemoteStatus = "merged"
)

type RemoteStorage interface {
	AssignReviewer(Reviewer) error
	GetStatus() (MrRemoteStatus, error)
}

type MrID string

type MergeRequest struct {
	id               MrID
	reviwe           *Rewive
	reviewer         *Reviewer
	state            MRState
	assignedReviewer *Reviewer
	createdAt        time.Time
	updatedAt        time.Time

	channel       NotificationChannel
	remoteStorage RemoteStorage
}
