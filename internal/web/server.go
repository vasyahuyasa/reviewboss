package web

import (
	"net/http"
)

type Server struct {
	//bot            *tgbotapi.BotAPI
	reviewHandlers *ReviewHandlers
	mux            *http.ServeMux
}

func NewServer(mux *http.ServeMux, reviewHandlers *ReviewHandlers) *Server {
	return &Server{
		//bot:            bot,
		reviewHandlers: reviewHandlers,
		mux:            mux,
	}
}

func (s *Server) RegisterRoutes() {
	s.mux.HandleFunc("/", s.reviewHandlers.Index)
	s.mux.HandleFunc("/register", s.reviewHandlers.Register)
	s.mux.HandleFunc("/reviwers", s.reviewHandlers.ListReviwers)
	/*
		mux.HandleFunc("/telegram/webhook", telegramHandler)
		_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://reviewboss.goshort.tk:443/telegram/webhook"))
		if err != nil {
			log.Fatal(err)
		}
	*/
}

func (s *Server) Listen(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

/*
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
*/
