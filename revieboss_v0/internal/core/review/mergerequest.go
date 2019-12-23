package review

import (
	"fmt"
)

const (
	StateNew state = iota
	StateWaitAccept
	StateWaitAcceptTimeout
	StateWaitProposal
	StateWaitProposalTimeout
	StateWaitProposalDeclined
	StateAccepted
)

// possible variants of state changes
var stateMap = map[state][]state{
	StateWaitAccept:           []state{StateAccepted, StateWaitAcceptTimeout},
	StateWaitAcceptTimeout:    []state{StateWaitProposal},
	StateWaitProposal:         []state{StateWaitProposalTimeout, StateAccepted, StateWaitProposalDeclined},
	StateWaitProposalTimeout:  []state{StateAccepted},
	StateWaitProposalDeclined: []state{StateWaitAccept},
	StateAccepted:             []state{},
}

type state int

// MergeRequest is piece of code for review
type MergeRequest struct {
	ID            string
	URL           string
	RequiredSkill SkillName
	currentState  state
}

func NewMergeRequest(ID string, URL string, RequiredSkill SkillName) *MergeRequest {
	return &MergeRequest{
		ID:            ID,
		URL:           URL,
		RequiredSkill: RequiredSkill,
		currentState:  StateNew,
	}
}

// Satate is return current state of merge request. Function is not thread safe
func (m *MergeRequest) Satate() state {
	return m.currentState
}

// SetState is change current state of mergerequest, if it is impossible change state or
// state not allowed error will be returned, if change state success after callback
// will be executed within lock.
func (m *MergeRequest) setState(newState state) error {
	states, ok := stateMap[newState]
	if !ok {
		return fmt.Errorf("unknow state %v", newState)
	}

	for _, possibleState := range states {
		if possibleState == newState {
			m.currentState = newState
			return nil
		}
	}

	return fmt.Errorf("state %v not allowed after state %v", newState, m.currentState)
}
