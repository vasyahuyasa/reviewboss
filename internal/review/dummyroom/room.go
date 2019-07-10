package dummyroom

import (
	"github.com/vasyahuyasa/reviewboss/internal/review"
)

type Room struct {
}

func (r *Room) List() ([]review.Reviewer, error) {
	return []review.Reviewer{
		{
			Name: "Vasya H.",
			Skills: []review.Skill{
				review.Skill{Name: "Go", Level: review.SkillHi},
				review.Skill{Name: "JavaScript", Level: review.SkillHi},
				review.Skill{Name: "Java", Level: review.SkillHi},
				review.Skill{Name: "C", Level: review.SkillHi},
			},
		},
		{
			Name: "Peter K.",
			Skills: []review.Skill{
				review.Skill{Name: "Go", Level: review.SkillHi},
				review.Skill{Name: "JavaScript", Level: review.SkillMiddle},
				review.Skill{Name: "PHP", Level: review.SkillLow},
			},
		},
		{
			Name: "Kirill S.",
			Skills: []review.Skill{
				review.Skill{Name: "Ruby", Level: review.SkillHi},
				review.Skill{Name: "Grioovy", Level: review.SkillMiddle},
			},
		},
	}, nil
}
