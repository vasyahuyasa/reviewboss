package domain

import (
	"fmt"
	"time"
)

// MRStateMachine is a concrete implementation
type mrStateMachine struct {
	mr *MergeRequest
}

func (sm *mrStateMachine) current() MRState {
	return sm.mr.state
}

// CanTransitionTo enforces allowed transitions
func (sm *mrStateMachine) CanTransitionTo(next MRState) bool {
	switch sm.current() {
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

func (sm *mrStateMachine) TransitionTo(next MRState) error {
	if !sm.CanTransitionTo(next) {
		return fmt.Errorf("invalid transition from %s to %s", sm.current(), next)
	}
	sm.mr.setState(next)
	sm.mr.setUpdatedAt(nowUnix())
	return nil
}

func nowUnix() time.Time {
	return time.Now()
}
