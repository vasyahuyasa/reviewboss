package review

import (
	"log"

	"github.com/pkg/errors"
)

type Service struct {
	Brain *Brain
}

func (s *Service) RegisterMergeRequest(id, url string, skill SkillName) (*MergeRequest, error) {
	s.Brain.Lock()
	defer s.Brain.Unlock()

	mr := NewMergeRequest(id, url, skill)
	err := s.Brain.RegisterMergeRequest(mr)
	if err != nil {
		return nil, errors.Wrapf(err, "can not register merge request %q", id)
	}

	reviwers, err := s.Brain.SelectReviwersBySkill(skill)
	if err != nil {
		return nil, errors.Wrapf(err, "can not get list of reviwers for skill %q", skill)
	}

	err = s.Brain.AssignReviwers(mr.ID, reviwers)
	if err != nil {
		return nil, errors.Wrapf(err, "can not assign reviwers to merge request %q", id)
	}

	return mr, nil
}

func (s *Service) ReviwersBySkill(skill SkillName) ([]Reviewer, error) {
	return s.Brain.SelectReviwersBySkill(skill)
}

func (s *Service) MergeRequestTimedOut(mr *MergeRequest) {
	s.Brain.Lock()
	defer s.Brain.Unlock()

	log.Printf("merge request %q timed out", mr.ID)

	ma, ok := s.Brain.findMergeAssign(mr.ID)
	if !ok {
		panic("merge assign not found but must be!")
	}

	err := ma.mergeRequest.setState(StateWaitProposal)
	if err != nil {
		log.Printf("can not set state of merge request %q to StateWaitProposal: ", err)
	}

	reviewer, ok := ma.firstReviwer()
	if !ok {
		log.Println("can not force set reviwer, list of reviwers is empty")
	}

	ma.setAssignee(reviewer)
	err = s.anonunceAssignee(ma.mergeRequest.ID, ma.mergeRequest.URL, reviewer.Name)
	if err != nil {
		log.Println("can not anounce assegnee: ", err)
	}
}

func (s *Service) anonunceAssignee(id, url, reviwer string) error {
	log.Printf("Reviwer %q has been assigned to merge request %q (%s)", reviwer, id, url)
	return nil
}
