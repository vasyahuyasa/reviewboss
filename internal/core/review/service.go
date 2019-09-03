package review

import "github.com/pkg/errors"

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
