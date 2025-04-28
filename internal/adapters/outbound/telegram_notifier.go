package outbound

import (
	"fmt"

	"github.com/vasyahuyasa/reviewboss/internal/domain"
	"github.com/vasyahuyasa/reviewboss/internal/ports"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramNotifier struct{ Bot *tgbotapi.BotAPI }

func (t *TelegramNotifier) NotifyProposal(mr *domain.MergeRequest, reviewer string) error {
	msg := fmt.Sprintf("Please review MR %s", mr.ID)
	_, err := t.Bot.Send(tgbotapi.NewMessage(int64(reviewer), msg))
	return err
}

// implement other methods...

var _ ports.Notifier = (*TelegramNotifier)(nil)
