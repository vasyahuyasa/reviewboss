package review

const (
	SkillLow SkillLevel = iota
	SkillMiddle
	SkillHi
)

// SkillLevel is describing reviwer level knowlege of some technology
type SkillLevel int

// Room is source of reviwers
type Room interface {
	List() ([]Reviewer, error)
}

// Skill is one of reviwer abilities
type Skill struct {
	Level SkillLevel
	Name  string
}

// Reviewer is person who should make revie
type Reviewer struct {
	Name   string
	Skills []Skill
}
