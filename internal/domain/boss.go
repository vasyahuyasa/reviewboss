package domain

import (
	"errors"
	"sync"
)

var (
	ErrMergeRequestAlreadyRegistered = errors.New("merge request already registred")
)

type Boss struct {
	mu            sync.Mutex
	mergeRequests []MergeRequest
}

func (b *Boss) AddMR(mr MergeRequest) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.isMrRegistered(mr) {
		return ErrMergeRequestAlreadyRegistered
	}

	return nil
}

func (b *Boss) ReviewerDeclined(mrid MrID, reviwer Reviewer) error {
	return nil
}

func (b *Boss) AssignReviwer(mrid MrID, reviwer Reviewer) error {
	return nil
}

func (b *Boss) isMrRegistered(mr MergeRequest) bool {
	return false
}

func (b *Boss) addMrToPool(mr MergeRequest) {
	b.mergeRequests = append(b.mergeRequests, mr)
}
