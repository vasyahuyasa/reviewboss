package web

import (
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Server struct {
	bot *tgbotapi.BotAPI
}

func NewServer(bot *tgbotapi.BotAPI) *Server {
	return &Server{
		bot: bot,
	}
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	/*
		mux.HandleFunc("/telegram/webhook", telegramHandler)
		_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://reviewboss.goshort.tk:443/telegram/webhook"))
		if err != nil {
			log.Fatal(err)
		}
	*/
}

func (s *Server) RunTelegram() {
	go s.waitForUpdates()
}

func (s *Server) waitForUpdates() {
	var lastMsg tgbotapi.Message
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := s.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("can not get update channel %v", err)
	}
	for update := range updates {
		if update.Message == nil {
			log.Printf("Unprocessible update: %v", update)
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		newMsg, err := s.bot.Send(msg)
		if err != nil {
			log.Fatalf("can not send message: %v", err)
		}

		if lastMsg.Chat == nil {
			continue
		}

		_, err = s.bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
			ChatID:    lastMsg.Chat.ID,
			MessageID: lastMsg.MessageID,
		})
		if err != nil {
			log.Printf("can not delete message %d: %v", lastMsg.MessageID, err)
		}
		lastMsg = newMsg
	}
}
