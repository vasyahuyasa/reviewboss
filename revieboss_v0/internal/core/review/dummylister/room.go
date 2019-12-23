package dummylister

import (
	"github.com/vasyahuyasa/reviewboss/internal/core/review"
)

type Lister struct {
}

func (r *Lister) List() ([]review.Reviewer, error) {
	return []review.Reviewer{
		{
			Name: "Vasya H.",
			Skills: []review.Skill{
				review.Skill{Name: "Go", Level: review.SkillLevelHi},
				review.Skill{Name: "JavaScript", Level: review.SkillLevelHi},
				review.Skill{Name: "Java", Level: review.SkillLevelHi},
				review.Skill{Name: "C", Level: review.SkillLevelHi},
			},
		},
		{
			Name: "Peter K.",
			Skills: []review.Skill{
				review.Skill{Name: "Go", Level: review.SkillLevelHi},
				review.Skill{Name: "JavaScript", Level: review.SkillLevelMiddle},
				review.Skill{Name: "PHP", Level: review.SkillLevelLow},
			},
		},
		{
			Name: "Ivan S.",
			Skills: []review.Skill{
				review.Skill{Name: "Ruby", Level: review.SkillLevelHi},
				review.Skill{Name: "PHP", Level: review.SkillLevelMiddle},
				review.Skill{Name: "Brainfuck", Level: review.SkillLevelLow},
			},
		},
	}, nil
}
