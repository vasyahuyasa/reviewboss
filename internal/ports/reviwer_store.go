package ports

type ReviewerStore interface {
	ListEligible(repoID, exclude string) ([]string, error)
	Assign(repoID, mrID, reviewerID string) error
}
