package review

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var ErrAlreadyRegistred = errors.New("merge request already registered")
var ErrNotFound = errors.New("merge request not found")

const (
	waitTimeout       = time.Minute
	waitAcceptTimeout = time.Minute
)

// TimoutHandler when timeout for merge request exceed
type TimoutHandler func(mr *MergeRequest)

// Brain is decede who should make review. Not thread-safe
type Brain struct {

	// TODO: dispose of global lock
	sync.Mutex

	reviwers ReviwerLister
	assigns  map[string]mergeAssign

	onTimeout         TimoutHandler
	onProposalTimeout TimoutHandler
}

type mergeAssign struct {
	mergeRequest *MergeRequest
	reviwers     *[]Reviewer
	assignee     Reviewer

	acceptCtx             context.Context
	acceptTimeoutCancel   context.CancelFunc
	proposalCtx           context.Context
	proposalTimeoutCancel context.CancelFunc
}

func NewBrain(reviwers ReviwerLister, onTimeout, onProposalTimeout TimoutHandler) *Brain {
	return &Brain{
		reviwers: reviwers,
		assigns:  map[string]mergeAssign{},

		onTimeout:         onTimeout,
		onProposalTimeout: onProposalTimeout,
	}
}

func (b *Brain) RegisterMergeRequest(mr *MergeRequest) error {
	_, ok := b.assigns[mr.ID]
	if ok {
		return ErrAlreadyRegistred
	}

	ctx, cancel := context.WithTimeout(context.Background(), waitAcceptTimeout)

	b.assigns[mr.ID] = mergeAssign{
		mergeRequest:        mr,
		acceptCtx:           ctx,
		acceptTimeoutCancel: cancel,
	}

	return nil
}

func (b *Brain) RemoveMergeRequest(id string) {
	delete(b.assigns, id)
}

func (b *Brain) SelectReviwers(requiredSkill SkillName) ([]Reviewer, error) {
	all, err := b.reviwers.List()
	if err != nil {
		return nil, errors.Wrap(err, "can not get list of reviwers")
	}

	// choose reviwers with requred skill and sort by skilllevel
	var suitable []Reviewer
	for _, r := range all {
		if r.SkillLevel(requiredSkill) != SkillLevelNone {
			suitable = append(suitable, r)
		}
	}

	// TODO: smart algo for sroting,  it should depend on review frequency
	sort.Slice(suitable, func(i, j int) bool {
		return suitable[i].SkillLevel(requiredSkill) > suitable[j].SkillLevel(requiredSkill)
	})

	return suitable, nil
}

func (b *Brain) AssignReviwers(id string, reviwers []Reviewer) error {
	_, ok := b.assigns[mr.ID]
	if !ok {
		return ErrNotFound
	}

}
