package main

import (
	"fmt"
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

type MergeRequestSource interface {
	ListOpenMergeRequests() ([]*MergeRequest, error)
	GetMergeRequest(id string) (*MergeRequest, error)
	SaveMergeRequest(mr *MergeRequest) error
}

type Reviewer struct {
	ID    string
	Name  string
	Email string
}

type MergeRequest struct {
	ID               string
	Title            string
	Author           string
	State            MRState
	AssignedReviewer *Reviewer
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type MergeRequestStateMachine struct {
	MR *MergeRequest
}

func (sm *MergeRequestStateMachine) TransitionTo(newState MRState) error {
	if !sm.isValidTransition(newState) {
		return fmt.Errorf("invalid transition from %s to %s", sm.MR.State, newState)
	}
	sm.MR.State = newState
	sm.MR.UpdatedAt = time.Now()
	// Save to DB here if needed
	return nil
}

func (sm *MergeRequestStateMachine) isValidTransition(newState MRState) bool {
	// You can hardcode rules here or load them dynamically
	switch sm.MR.State {
	case StateNew:
		return newState == StateWaitingVoluntary
	case StateWaitingVoluntary:
		return newState == StateAgreementReceived || newState == StateVoluntaryExpired
	case StateVoluntaryExpired:
		return newState == StateProposeReviewer
	case StateProposeReviewer:
		return newState == StateReviewerDeclined || newState == StateAgreementReceived || newState == StateProposalExpired
	case StateReviewerDeclined:
		return newState == StateProposeReviewer || newState == StateNoReviewers
	case StateNoReviewers:
		return newState == StateWaitingSomeone
	case StateWaitingSomeone:
		return newState == StateAgreementReceived || newState == StateWaitingSomeoneExpired
	case StateWaitingSomeoneExpired:
		return newState == StateAssignReviewer
	case StateAgreementReceived, StateAssignReviewer:
		return newState == StateReviewInProgress
	case StateReviewInProgress:
		return newState == StateMerged || newState == StateLongInReview
	case StateLongInReview:
		return newState == StateMerged
	default:
		return false
	}
}
