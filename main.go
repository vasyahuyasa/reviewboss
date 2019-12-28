package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	gitlabapi "github.com/xanzy/go-gitlab"
	"golang.org/x/net/proxy"
)

const (
	statusNone mergeRequestStatus = iota
	statusNew
	statusWaitAssignee
	statusAssigned
	statusDone

	timeoutNew          = time.Minute
	timeoutWaitAseegnee = time.Minute
	timeoutAssigneed    = time.Minute

	configFile = "config.toml"
)

var (

	// map[telegramID][gitlabLogin]
	telegram2gitlab = map[int]string{
		// TODO: machanic of register account
		// TODO: cahe of names for gitlab projects
		421340245: "ioaioa",
		310975055: "ioaioa",
	}
)

type config struct {
	BotToken    string `toml:"bot_token"`
	Proxy       string `toml:"proxy"`
	ReMatch     string `toml:"re_match"`
	GitlabToken string `toml:"gitlab_token"`
}

type reviwer struct {
	GitlabID     int
	TelegramID   int
	TelegramName string
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
	FromID           int

	Link                  string
	GitlabProject         string
	GitlabMergeRequesetID int
}

type mergeRequestInput struct {
	link                 string
	gitlabProject        string
	gitlabMergeRequestId int
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

type gitlab struct {
	api *gitlabapi.Client
}

func (input *mergeRequestInput) ID() string {
	return fmt.Sprintf("%s_%d", input.gitlabProject, input.gitlabMergeRequestId)
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
							TelegramName: "reviwer1",
						},
						{
							TelegramName: "reviwer2",
						},
					}

					eng.mergeRequests[id] = m
					if len(m.ProposeAssignees) > 0 {
						assignees := ""
						for i, proposedAssignee := range m.ProposeAssignees {
							if i != 0 {
								assignees += ", " + proposedAssignee.TelegramName
							} else {
								assignees += proposedAssignee.TelegramName
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
						r := m.ProposeAssignees[0]
						m.Assignee = &r

						log.Printf("%q moved to status statusAssigned with reviwer %q", m.ID, m.Assignee.TelegramName)

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
		return errors.New("shutdown timed out, force quit")
	case <-eng.done:
	}

	return nil
}

func (git *gitlab) Assign(project string, mergeRequest int, assigneeID int) error {
	_, _, err := git.api.MergeRequests.UpdateMergeRequest(project, mergeRequest, &gitlabapi.UpdateMergeRequestOptions{
		AssigneeIDs: []int{assigneeID},
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

func extractMergeLinks(s string, re *regexp.Regexp) []mergeRequestInput {
	var input []mergeRequestInput
	for _, match := range re.FindAllStringSubmatch(s, -1) {
		id, err := strconv.Atoi(match[2])
		if err != nil {
			log.Printf("can not convert %q to int: %v", match[2], err)
			continue
		}

		input = append(input, mergeRequestInput{
			link:                 match[0],
			gitlabProject:        match[1],
			gitlabMergeRequestId: id,
		})
	}

	return input
}

func telegramToGitlab(git *gitlab, telegramUserID int, gitlabProject string) (int, error) {
	username, ok := telegram2gitlab[telegramUserID]
	if !ok {
		return 0, errors.New("user not registed")
	}

	userID, err := git.ProjectUserID(gitlabProject, username)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func main() {
	var cfg config
	_, err := toml.DecodeFile(configFile, &cfg)
	if err != nil {
		log.Fatalf("can not read config %q: %v", configFile, err)
	}

	// regular expression for extract gitlab links
	reMergeRequest := regexp.MustCompile(cfg.ReMatch)

	// telegram bot
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

	bot, err := tgbotapi.NewBotAPIWithClient(cfg.BotToken, httpClient)
	if err != nil {
		log.Fatalf("can not inittialize telegram bot: %v", err)
	}

	log.Printf("telegram bot authorized on account %s", bot.Self.UserName)

	// gitlab api
	git := &gitlab{
		api: gitlabapi.NewClient(&http.Client{}, cfg.GitlabToken),
	}

	// main logic
	eng := engine{
		running:       true,
		done:          make(chan interface{}),
		mergeRequests: map[string]mergeRequest{},

		onStatusWaitAssignee: func(m mergeRequest) {
			text := fmt.Sprintf("[%s](%s)", m.Link, m.Link)
			for i, r := range m.ProposeAssignees {
				if i == 0 {
					text += "\n\n*Кандидаты на проведение ревью:*"
				}
				text += fmt.Sprintf("\n%d. %s", i+1, r.TelegramName)
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
		u := tgbotapi.UpdateConfig{Timeout: 60}

		updates, err := bot.GetUpdatesChan(u)
		if err != nil {
			log.Fatalf("telegram can not get update channel: %v", err)
		}
		time.Sleep(time.Millisecond * 500)
		updates.Clear()

		for update := range updates {
			// text message on channel
			if update.Message != nil {
				log.Println("From", update.Message.From.String(), update.Message.From.ID)

				for _, input := range extractMergeLinks(update.Message.Text, reMergeRequest) {
					log.Println("MergeRequest:", input.gitlabProject, input.gitlabMergeRequestId)
					err := eng.AddMergeRequest(mergeRequest{
						ID:        input.ID(),
						Status:    statusNew,
						AddedOn:   time.Now(),
						UpdatedOn: time.Now(),
						ChatID:    update.Message.Chat.ID,
						From:      update.Message.From.String(),
						FromID:    update.Message.From.ID,

						Link:                  input.link,
						GitlabProject:         input.gitlabProject,
						GitlabMergeRequesetID: input.gitlabMergeRequestId,
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
					// TODO: telegram ID to gitlab id
					m.Assignee = &reviwer{
						TelegramName: update.CallbackQuery.From.String(),
						TelegramID:   update.CallbackQuery.From.ID,
					}
					eng.mergeRequests[id] = m

					err := git.Assign(m.GitlabProject, m.GitlabMergeRequesetID, m.Assignee.GitlabID)
					if err != nil {
						log.Printf("can not assign %q to merge request %q: %v", m.Assignee.TelegramName, m.ID, err)
					}

					projectUserID, err := telegramToGitlab(git, update.CallbackQuery.From.ID, m.GitlabProject)
					if err != nil {
						log.Printf("can not get userID of %q for project %q: %v", update.CallbackQuery.From.String(), m.GitlabProject, err)
					} else {
						log.Println("gitlab user:", projectUserID)
					}
				}

				eng.mu.Unlock()

				/*
					newKbd := tgbotapi.NewEditMessageReplyMarkup(
						update.CallbackQuery.Message.Chat.ID,
						update.CallbackQuery.Message.MessageID,
						tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{}),
					)
				*/

				newMsg := tgbotapi.NewEditMessageText(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
					fmt.Sprintf("[%s](%s)\n\n*Ревью проводит:* @%s", m.Link, m.Link, m.Assignee.TelegramName),
				)
				newMsg.ParseMode = "markdown"

				_, err := bot.Send(newMsg)
				if err != nil {
					log.Printf("can not edit keyboard: %v", err)
				}
			}
		}
	}()

	time.Sleep(time.Minute * 10)
	err = eng.Shutdown()
	if err != nil {
		log.Printf("can not shutdown engine: %v", err)
	}
}
