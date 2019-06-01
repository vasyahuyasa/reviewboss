package web

import (
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Server struct {
	bot *tgbotapi.Bot
}

func NewServer(bot *tgbotapi.Bot) *Server {
	return &Server{
		bot: bot,
	}
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/telegram/webhook", telegramHandler)
	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://reviewboss.goshort.tk:443/telegram/webhook"))
	if err != nil {
		log.Fatal(err)
	}
}
