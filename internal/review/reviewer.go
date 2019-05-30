package review

type Telegram struct {
	Username string
}

type Gitlab struct {
}

type Reviewer struct {
	Skills   []Skill
	Telegram Telegram
	Gitlab   Gitlab
}
