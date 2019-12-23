package review

import (
	"log"
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
	assigns  map[string]*mergeAssign

	onTimeout         TimoutHandler
	onProposalTimeout TimoutHandler
}

type mergeAssign struct {
	mergeRequest *MergeRequest
	reviwers     []Reviewer
	assignee     Reviewer
	hasAssignee  bool

	acceptTimeoutCancel   chan struct{}
	proposalTimeoutCancel chan struct{}
}

func NewBrain(reviwers ReviwerLister, onTimeout, onProposalTimeout TimoutHandler) *Brain {
	return &Brain{
		reviwers: reviwers,
		assigns:  map[string]*mergeAssign{},

		onTimeout:         onTimeout,
		onProposalTimeout: onProposalTimeout,
	}
}

func (b *Brain) RegisterMergeRequest(mr *MergeRequest) error {
	_, ok := b.findMergeAssign(mr.ID)
	if ok {
		return ErrAlreadyRegistred
	}

	cancel := make(chan struct{})
	go func() {
		select {
		case <-time.After(waitAcceptTimeout):
			err := mr.setState(StateWaitAcceptTimeout)
			if err != nil {
				log.Println("get waitAceptTimeout timer but can not set new state: ", err)
				return
			}
			b.onTimeout(mr)

		case <-cancel:
		}
	}()

	err := mr.setState(StateWaitAccept)
	if err != nil {
		log.Println("can not set state of merge request to StateWaitAccept: ", err)
	}

	b.assigns[mr.ID] = &mergeAssign{
		mergeRequest:        mr,
		acceptTimeoutCancel: cancel,
	}

	return nil
}

func (b *Brain) RemoveMergeRequest(id string) {
	delete(b.assigns, id)
}

func (b *Brain) SelectReviwersBySkill(requiredSkill SkillName) ([]Reviewer, error) {
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
	mergeRequest, ok := b.assigns[id]
	if !ok {
		return ErrNotFound
	}

	mergeRequest.reviwers = reviwers

	return nil
}

func (b *Brain) findMergeAssign(id string) (*mergeAssign, bool) {
	m, ok := b.assigns[id]
	return m, ok
}

func (ma *mergeAssign) firstReviwer() (Reviewer, bool) {
	if len(ma.reviwers) == 0 {
		return Reviewer{}, false
	}

	return ma.reviwers[0], true
}

func (ma *mergeAssign) setAssignee(reviwer Reviewer) {
	ma.assignee = reviwer
	ma.hasAssignee = true
}

func (ma *mergeAssign) removeAssignee() {
	ma.assignee = Reviewer{}
	ma.hasAssignee = false
}
