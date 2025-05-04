package domain

type mergeRequestCollection struct {
	mergeRequests []MergeRequest
}

func (c *mergeRequestCollection) add(mr MergeRequest) {
	c.mergeRequests = append(c.mergeRequests, mr)
}

func (c *mergeRequestCollection) get(id MrID) (MergeRequest, bool) {
	for _, mr := range c.mergeRequests {
		if mr.getID() == id {
			return mr, true
		}
	}

	return MergeRequest{}, false
}
