package review

import (
	"sort"
	"strings"

	"github.com/vasyahuyasa/reviewboss/internal/infrasturcure/notify"

	"github.com/pkg/errors"
)

var ErrAlreadyRegistred = errors.New("merge request already registered")

// MergeRequest is piece of code for review
type MergeRequest struct {
	ID            string
	URL           string
	RequiredSkill SkillName
}

type mergeAssign struct {
	merge    MergeRequest
	reviwers []Reviewer
	assignee *Reviewer
}

// Brain is decede who should make review
type Brain struct {
	reviwers ReviwerLister
	assigns  map[string]mergeAssign
}

func NewBrain(reviwers ReviwerLister) *Brain {
	return &Brain{
		reviwers: reviwers,
		assigns:  map[string]mergeAssign{},
	}
}

func (b *Brain) Register(mr MergeRequest) error {
	_, ok := b.assigns[mr.ID]
	if ok {
		return ErrAlreadyRegistred
	}

	b.assigns[mr.ID] = mergeAssign{
		merge: mr,
	}

	reviwers, err := b.SelectReviwers(mr.RequiredSkill)
	if err != nil {
		return err
	}

	list := []string{}
	for _, r := range reviwers {
		list = append(list, r.Name)
	}

	notify.ShowMessage("Список ревьюверов: " + strings.Join(list, ", "))

	return nil
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
