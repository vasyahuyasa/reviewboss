package inbound

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/vasyahuyasa/reviewboss/internal/domain"
	"github.com/vasyahuyasa/reviewboss/internal/ports"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramSource listens for MR links in Telegram chats
// and creates new MergeRequest entities for processing.
type TelegramSource struct {
	Bot            *tgbotapi.BotAPI
	Store          ports.MergeRequestSource
	StateMachineFn func(mr *domain.MergeRequest) domain.StateMachine
	ReviewService  ports.ReviewerService
}

// Run starts polling Telegram updates
func (t *TelegramSource) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.Bot.GetUpdatesChan(u)

	// regex to match MR link, e.g. https://gitlab.com/org/repo/-/merge_requests/123
	re := regexp.MustCompile(`https?://[^/]+/([^/]+/[^/]+)/(?:-/)?merge_requests/(\d+)`)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			text := update.Message.Text
			match := re.FindStringSubmatch(text)
			if len(match) != 3 {
				continue // no MR link
			}
			repoID := match[1]
			mrID := match[2]

			// initialize MR in domain
			mr := &domain.MergeRequest{
				ID:        mrID,
				RepoID:    repoID,
				Author:    strings.TrimSpace(update.Message.From.UserName),
				State:     domain.StateNew,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			}
			// persist MR
			if err := t.Store.Save(mr); err != nil {
				// optionally notify error
				continue
			}
			// bootstrap state machine and transition to waiting
			sm := t.StateMachineFn(mr)
			sm.TransitionTo(domain.StateWaitingVoluntary)

			// propose reviewer on Telegram
			t.ReviewService.Propose(mr)
		}
	}
}

var _ ports.MRSource = (*TelegramSource)(nil)
