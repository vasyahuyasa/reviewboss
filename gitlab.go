package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	gitlabapi "github.com/xanzy/go-gitlab"
)

const (
	mergeStateOpened = "opened"
	mergeStateClosed = "closed"
	mergeStateLocked = "locked"
	mergeStateMerged = "merged"
)

type gitlabUser struct {
	ID       int
	Username string
}

type gitlabMergeRequest struct {
	ID       int
	Title    string
	State    string
	Changes  int
	MergedAt time.Time
	ClosedAt time.Time
	WebURL   string
	Labels   []string
	Author   *gitlabUser
	Assignee *gitlabUser
}

type gitlab struct {
	api *gitlabapi.Client
}

func (git *gitlab) Assign(project string, mergeRequest int, assigneeID int) error {
	_, _, err := git.api.MergeRequests.UpdateMergeRequest(project, mergeRequest, &gitlabapi.UpdateMergeRequestOptions{
		AssigneeID: &assigneeID,
	})

	return err
}

func (git *gitlab) ProjectUserID(project, username string) (int, error) {
	users, _, err := git.api.Projects.ListProjectsUsers(project, &gitlabapi.ListProjectUserOptions{
		Search: &username,
	})
	if err != nil {
		return 0, err
	}

	if len(users) == 0 {
		return 0, errors.New("user have no access to project")
	}

	return users[0].ID, nil
}

// MergeRequestInfo return title and changed file count for merge request
func (git *gitlab) MergeRequestInfo(project string, mergeRequest int) (gitlabMergeRequest, error) {
	mr, _, err := git.api.MergeRequests.GetMergeRequestChanges(project, mergeRequest)
	if err != nil {
		return gitlabMergeRequest{}, err
	}

	cnt, err := strconv.Atoi(mr.ChangesCount)
	if err != nil {
		return gitlabMergeRequest{}, fmt.Errorf("can not parse number of changes for merge request: %w", err)
	}

	var author *gitlabUser
	var assignee *gitlabUser

	if mr.Author != nil {
		author = &gitlabUser{
			ID:       mr.Author.ID,
			Username: mr.Author.Username,
		}
	}

	if mr.Assignee != nil {
		assignee = &gitlabUser{
			ID:       mr.Assignee.ID,
			Username: mr.Assignee.Username,
		}
	}

	return gitlabMergeRequest{
		ID:       mr.ID,
		Title:    mr.Title,
		State:    mr.State,
		Changes:  cnt,
		WebURL:   mr.WebURL,
		Labels:   mr.Labels,
		Author:   author,
		Assignee: assignee,
	}, nil
}
