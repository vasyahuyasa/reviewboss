package review

import (
	"errors"
)

var ErrAlreadyRegistred = errors.New("merge request already registered")

// Brain is decede who should make review
type Brain struct {
	room  Room
	mreqs map[int]MergeRequest
}

func NewBrain(room Room) *Brain {
	return &Brain{
		room: room,
	}
}

func (b *Brain) Register(mr MergeRequest) error {
	_, ok := b.mreqs[mr.ID]
	if ok {
		return ErrAlreadyRegistred
	}

	b.mreqs[mr.ID] = mr
	return nil
}
