package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/proxy"
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

var reMergeRequest = regexp.MustCompile(`(?Um)https://gitlab\.devpizzasoft\.ru:8000/(.+)/merge_requests/(\d+)\b`)

type config struct {
	Token         string `toml:"token"`
	Proxy         string `toml:"proxy"`
	MergeRequstRe string `toml:"re_match"`
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
	ChatID           int64
	From             string
	Link             string
}

type mergeRequestInput struct {
	link    string
	project string
	id      string
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

func (input *mergeRequestInput) ID() string {
	return input.project + "_" + input.id
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

func extractMergeLinks(s string) []mergeRequestInput {
	// https://gitlab.devpizzasoft.ru:8000/pizzasoft/socket/merge_requests/75
	var input []mergeRequestInput
	for _, match := range reMergeRequest.FindAllStringSubmatch(s, -1) {
		input = append(input, mergeRequestInput{
			link:    match[0],
			project: match[1],
			id:      match[2],
		})
	}

	return input
}

func main() {
	var cfg config
	_, err := toml.DecodeFile(configFile, &cfg)
	if err != nil {
		log.Fatalf("can not read config %q: %v", configFile, err)
	}

	var httpClient *http.Client

	if cfg.Proxy != "" {
		dialer, err := proxy.SOCKS5("tcp", cfg.Proxy, nil, nil)
		if err != nil {
			log.Fatalf("can't connect to the proxy: %v", err)
		}

		httpClient = &http.Client{Transport: &http.Transport{Dial: dialer.Dial}}
	} else {
		httpClient = &http.Client{}
	}

	bot, err := tgbotapi.NewBotAPIWithClient(cfg.Token, httpClient)
	if err != nil {
		log.Fatalf("can not inittialize telegram bot: %v", err)
	}

	log.Printf("telegram bot authorized on account %s", bot.Self.UserName)

	eng := engine{
		running:       true,
		done:          make(chan interface{}),
		mergeRequests: map[string]mergeRequest{},

		onStatusWaitAssignee: func(m mergeRequest) {
			text := fmt.Sprintf("[%s](%s) ожидает ревью", m.Link, m.Link)
			for i, r := range m.ProposeAssignees {
				if i == 0 {
					text += "\n\nКандидаты на проведение ревью:\n\n"
				}
				text += fmt.Sprintf("\n%d. %s", i+1, r.ID)
			}

			msg := tgbotapi.NewMessage(m.ChatID, text)
			msg.ParseMode = "markdown"
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Я возьму", m.ID),
			})
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("can not send message to channel %d: %v", m.ChatID, err)
			}
		},
		onStatusAssigned:      func(m mergeRequest) {},
		onNoProposedAssignees: func(m mergeRequest) {},
		onDone:                func(m mergeRequest) {},
		beforeClean:           func(m mergeRequest) {},
	}

	go eng.Watcher()

	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, err := bot.GetUpdatesChan(u)
		if err != nil {
			log.Fatalf("telegram can not get update channel: %v", err)
		}
		time.Sleep(time.Millisecond * 500)
		updates.Clear()

		for update := range updates {

			log.Printf("%v", update)

			// text message on channel
			if update.Message != nil {
				for _, input := range extractMergeLinks(update.Message.Text) {
					log.Println("MergeRequest:", input.project, input.id)
					err := eng.AddMergeRequest(mergeRequest{
						ID:        input.ID(),
						Status:    statusNew,
						AddedOn:   time.Now(),
						UpdatedOn: time.Now(),
						ChatID:    update.Message.Chat.ID,
						From:      update.Message.From.String(),
						Link:      input.link,
					})

					if err != nil {
						log.Printf("can not register mergerequest %q: %v", input.ID(), err)
					}
				}
			}

			// somebody agreed to take review
			if update.CallbackQuery != nil {
				id := update.CallbackQuery.Data
				eng.mu.Lock()
				m, ok := eng.mergeRequests[id]

				if !ok {
					eng.mu.Unlock()
					continue
				}

				if m.Assignee == nil {
					m.Assignee = &reviwer{ID: update.CallbackQuery.From.String()}
					eng.mergeRequests[id] = m
				}

				eng.mu.Unlock()
				/*
					tgbotapi.NewEditMessageReplyMarkup(
						update.CallbackQuery.Message.Chat.ID,
						update.CallbackQuery.Message.MessageID)
				*/
				log.Println("callback", update.CallbackQuery)
			}
		}
	}()

	time.Sleep(time.Minute * 10)
	err = eng.Shutdown()
	if err != nil {
		log.Printf("can not shutdown engine: %v", err)
	}
}
