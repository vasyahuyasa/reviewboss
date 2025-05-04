package domain

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrMergeRequestAlreadyRegistered = errors.New("merge request already registered")
)

type Boss struct {
	mu            sync.Mutex
	mergeRequests mergeRequestCollection
}

func (b *Boss) AddMR(id MrID, channel NotificationChannel, storage RemoteStorage) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.isMrRegistered(id) {
		return ErrMergeRequestAlreadyRegistered
	}

	mr := MergeRequest{
		id:            id,
		state:         StateNew,
		createdAt:     b.now(),
		updatedAt:     b.now(),
		channel:       channel,
		remoteStorage: storage,
	}

	b.addMrToPool(mr)

	return nil
}

func (b *Boss) ReviewerDeclined(id MrID) error {
	return nil
}

func (b *Boss) AssignReviewer(id MrID, reviewer Reviewer) error {
	return nil
}

func (b *Boss) isMrRegistered(id MrID) bool {
	_, ok := b.mergeRequests.get(id)
	return ok
}

func (b *Boss) addMrToPool(mr MergeRequest) {
	b.mergeRequests.add(mr)
}

func (b *Boss) now() time.Time {
	return time.Now()
}
