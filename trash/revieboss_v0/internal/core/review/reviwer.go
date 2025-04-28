package review

const (
	SkillLevelNone SkillLevel = iota
	SkillLevelLow
	SkillLevelMiddle
	SkillLevelHi
)

const (
	SkillGolang SkillName = "Go"
)

// SkillLevel is describing reviwer level knowlege of some technology
type SkillLevel int

// ReviwerLister is source of reviwers
type ReviwerLister interface {
	List() ([]Reviewer, error)
}

// SkillName is name of tecnology
type SkillName string

// Skill is one of reviwer abilities
type Skill struct {
	Level SkillLevel
	Name  SkillName
}

// Reviewer is person who should make revie
type Reviewer struct {
	Name   string
	Skills []Skill
}

// SkillLevel representation of reviwer ability to selected skill
func (r Reviewer) SkillLevel(skill SkillName) SkillLevel {
	for _, s := range r.Skills {
		if s.Name == skill {
			return s.Level
		}
	}
	return SkillLevelNone
}

func (r Reviewer) TotalSkill() SkillLevel {
	var total SkillLevel
	for _, s := range r.Skills {
		total = total + s.Level
	}
	return total
}
