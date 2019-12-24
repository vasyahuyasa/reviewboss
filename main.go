package main

import (
	"errors"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

	configFile = "config.toml"
)

type config struct {
	Token   string `toml:"token"`
	Proxy   string `toml:"proxy"`
	GroupID int64  `toml:"groupid"`
}

type reviwer struct {
	ID string
}

type mergeRequestStatus int

type mergeRequest struct {
	ID               string
	Status           mergeRequestStatus
	AddedOn          time.Time
	UpdatedOn        time.Time
	Assignee         *reviwer
	ProposeAssignees []reviwer
}

type engine struct {
	running       bool
	mu            sync.RWMutex
	mergeRequests map[string]mergeRequest
	done          chan interface{}

	// new -> waitAssignee
	onStatusWaitAssignee func(m mergeRequest)

	// waitAssegnee -> assigned
	onStatusAssigned func(m mergeRequest)

	// waitAssegnee -> waitAssegnee because no on to assign
	onNoProposedAssignees func(m mergeRequest)

	// assigned -> done
	onDone func(m mergeRequest)

	// before clean done
	beforeClean func(m mergeRequest)
}

func (m *mergeRequest) SetStatus(status mergeRequestStatus) {
	m.Status = status
	m.UpdatedOn = time.Now()
}

func (eng *engine) AddMergeRequest(m mergeRequest) error {
	eng.mergeRequests[m.ID] = m
	return nil
}

func (eng *engine) Watcher() {
	for eng.running {
		eng.mu.Lock()
		for id, m := range eng.mergeRequests {
			timeout := time.Since(m.UpdatedOn)

			switch m.Status {
			case statusNew:
				if timeout >= timeoutNew {

					// logging status change
					m.SetStatus(statusWaitAssignee)

					// TODO: select list of assignees

					m.ProposeAssignees = []reviwer{
						{
							ID: "reviwer1",
						},
						{
							ID: "reviwer1",
						},
					}

					eng.mergeRequests[id] = m
					if len(m.ProposeAssignees) > 0 {
						assignees := ""
						for i, proposedAssignee := range m.ProposeAssignees {
							if i != 0 {
								assignees += ", " + proposedAssignee.ID
							} else {
								assignees += proposedAssignee.ID
							}

						}
						log.Printf("%q moved to status statusWaitAssignee with proposed assignees %s", m.ID, assignees)
					} else {
						log.Printf("%q moved to status statusWaitAssignee without proposed assignees", m.ID)
					}

					// TODO: message with list of aseegnees
					eng.onStatusWaitAssignee(m)
				}
			case statusWaitAssignee:
				if timeout >= timeoutWaitAseegnee {
					if len(m.ProposeAssignees) > 0 {
						m.SetStatus(statusAssigned)

						// TODO: force aseegnee someone
						m.Assignee = &reviwer{ID: m.ProposeAssignees[0].ID}

						log.Printf("%q moved to status statusAssigned with reviwer %q", m.ID, m.Assignee.ID)

						eng.onStatusAssigned(m)
					} else {
						m.SetStatus(statusWaitAssignee)
						log.Printf("%q keep in status statusWaitAssignee because merge request don't have proposed assignes", m.ID)
						eng.onNoProposedAssignees(m)
					}

					eng.mergeRequests[id] = m
				}
			case statusAssigned:
				if timeout >= timeoutAssigneed {
					// TODO: check merge request status and maybe set done
					// if mr.Done {
					//     m.SetStatus(statusDone)
					// }
					if false {
						m.SetStatus(statusDone)
						eng.mergeRequests[id] = m
						eng.onDone(m)
					}
				}
			case statusDone:
				eng.beforeClean(m)
				// TODO: remove mergerequest
			default:
				log.Printf("merge request %q in statusNone state, removing", m.ID)
				// TODO: remove mergerequest
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
		return errors.New("timeout")
	case <-eng.done:
	}

	return nil
}

func main() {
	var cfg config
	_, err := toml.DecodeFile(configFile, &cfg)
	if err != nil {
		log.Fatalf("can not read config %q: %v", configFile, err)
	}

	bot := botAPI{
		token:     cfg.Token,
		proxyAddr: cfg.Proxy,
		onUpdate: func(update tgbotapi.Update) {
			if update.Message == nil || update.Message.Chat == nil {
				return
			}

			// process only links
			u, err := url.Parse(update.Message.Text)
			if err != nil {
				return
			}

			log.Println("url:", u)
		},
	}

	err = bot.waitForUpdates()
	if err != nil {
		log.Fatalf("telegram bot is off: %v", err)
	}

	eng := engine{
		running:       true,
		done:          make(chan interface{}),
		mergeRequests: map[string]mergeRequest{},

		onStatusWaitAssignee: func(m mergeRequest) {},
		onStatusAssigned: func(m mergeRequest) {
			//msg := tgbotapi.NewMessage(cfg.GroupID, fmt.Sprintf("%q is assigned to %s", m.ID, "@ioaioa"))
			// _, err := bot.Send(msg)
			// if err != nil {
			// 	log.Printf("can not send to telegram: %v", err)
			// }
		},
		onNoProposedAssignees: func(m mergeRequest) {},
		onDone:                func(m mergeRequest) {},
		beforeClean:           func(m mergeRequest) {},
	}

	go eng.Watcher()

	eng.AddMergeRequest(mergeRequest{
		ID:        "mr1",
		Status:    statusNew,
		AddedOn:   time.Now(),
		UpdatedOn: time.Now(),
	})

	eng.AddMergeRequest(mergeRequest{
		ID:        "mr2",
		Status:    statusNew,
		AddedOn:   time.Now(),
		UpdatedOn: time.Now(),
	})

	time.Sleep(time.Second * 100)
	log.Println(eng.mergeRequests)
	err = eng.Shutdown()
	if err != nil {
		log.Printf("can not shutdown engine: %v", err)
	}
}
