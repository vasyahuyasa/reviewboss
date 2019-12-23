package main

import (
	"errors"
	"log"
	"sync"
	"time"
)

const (
	statusNone mergeRequestStatus = iota
	statusNew
	statusWaitAssignee
	statusAssigned
	statusDone

	timeoutNew          = time.Second
	timeoutWaitAseegnee = time.Second * 2
	timeoutAssigneed    = time.Second * 5
)

type reviwer struct {
	ID string
}

type mergeRequestStatus int

type mergeRequest struct {
	ID        string
	Status    mergeRequestStatus
	AddedOn   time.Time
	UpdatedOn time.Time
	Assignee  *reviwer
}

type engine struct {
	running       bool
	mu            sync.RWMutex
	mergeRequests map[string]mergeRequest
	done          chan interface{}
}

func (eng *engine) AddMergeRequest(m mergeRequest) error {
	eng.mergeRequests[m.ID] = m
	return nil
}

func (eng *engine) SetAseegnee(mergeID string, aseegnee *reviwer) error {
	eng.mu.Lock()
	defer eng.mu.Unlock()

	m, ok := eng.mergeRequests[mergeID]
	if !ok {
		return errors.New("merge request not found")
	}

	m.Assignee = aseegnee
	eng.mergeRequests[mergeID] = m
	return nil
}

func (eng *engine) Watcher() {
	for eng.running {
		eng.mu.Lock()
		for _, m := range eng.mergeRequests {
			timeout := time.Since(m.UpdatedOn)

			switch m.Status {
			case statusNew:
				if timeout >= timeoutNew {
					// TODO: message with list of aseegnees
					m.Status = statusWaitAssignee
				}
			case statusWaitAssignee:
				if timeout >= timeoutWaitAseegnee {
					// TODO: force aseegnee someone
					m.Status = statusAssigned
				}
			case statusAssigned:
				if timeout >= timeoutAssigneed {
					// TODO: check merge request status and maybe set done
				}
			case statusDone:
			}
		}
		eng.mu.Unlock()

		time.Sleep(time.Second)
	}

	close(eng.done)
}

func (eng *engine) Shutdown() error {
	eng.running = false

	select {
	case <-time.After(time.Second * 10):
	case <-eng.done:
	}

	return nil
}

func main() {
	eng := engine{
		running:       true,
		done:          make(chan interface{}),
		mergeRequests: map[string]mergeRequest{},
	}

	go eng.Watcher()

	eng.AddMergeRequest(mergeRequest{
		ID:        "mr1",
		Status:    statusNew,
		AddedOn:   time.Now(),
		UpdatedOn: time.Now(),
	})

	log.Println(eng.mergeRequests)
	eng.Shutdown()
}
