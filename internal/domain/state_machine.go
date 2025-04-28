package domain

import (
	"fmt"
	"time"
)

// StateMachine defines operations for MR state transitions
type StateMachine interface {
	Current() MRState
	CanTransitionTo(MRState) bool
	TransitionTo(MRState) error
}

// MRStateMachine is a concrete implementation
type MRStateMachine struct {
	MR *MergeRequest
}

func (sm *MRStateMachine) Current() MRState {
	return sm.MR.State
}

// CanTransitionTo enforces allowed transitions
func (sm *MRStateMachine) CanTransitionTo(next MRState) bool {
	switch sm.MR.State {
	case StateNew:
		return next == StateWaitingVoluntary
	case StateWaitingVoluntary:
		return next == StateAgreementReceived || next == StateVoluntaryExpired
	case StateVoluntaryExpired:
		return next == StateProposeReviewer
	case StateProposeReviewer:
		return next == StateReviewerDeclined || next == StateAgreementReceived || next == StateProposalExpired
	case StateReviewerDeclined:
		return next == StateProposeReviewer || next == StateNoReviewers
	case StateNoReviewers:
		return next == StateWaitingSomeone
	case StateWaitingSomeone:
		return next == StateAgreementReceived || next == StateWaitingSomeoneExpired
	case StateWaitingSomeoneExpired:
		return next == StateAssignReviewer
	case StateAgreementReceived, StateAssignReviewer:
		return next == StateReviewInProgress
	case StateReviewInProgress:
		return next == StateMerged || next == StateLongInReview
	case StateLongInReview:
		return next == StateMerged
	default:
		return false
	}
}

func (sm *MRStateMachine) TransitionTo(next MRState) error {
	if !sm.CanTransitionTo(next) {
		return fmt.Errorf("invalid transition from %s to %s", sm.MR.State, next)
	}
	sm.MR.State = next
	sm.MR.UpdatedAt = nowUnix()
	return nil
}

func nowUnix() time.Time {
	return time.Now()
}
