package main

import (
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
	Asssignee *reviwer
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

func (eng *engine) Watcher() {
	for eng.running {
		eng.mu.Lock()
		for _, m := range eng.mergeRequests {
			timeout := time.Since(m.UpdatedOn)

			switch m.Status {
			case statusNew:
				if timeout >= timeoutNew {

				}
			case statusWaitAssignee:
				if timeout >= timeoutWaitAseegnee {

				}
			case statusAssigned:
				if timeout >= timeoutAssigneed {

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
