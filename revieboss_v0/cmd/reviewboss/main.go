package main

import (
	"log"
	"net/http"

	"github.com/vasyahuyasa/reviewboss/internal/web"

	"database/sql"

	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vasyahuyasa/reviewboss/internal/core/review"
	"github.com/vasyahuyasa/reviewboss/internal/core/review/dummylister"
	gitlab "github.com/xanzy/go-gitlab"
)

func main() {
	cfg := loadConfig()

	git := gitlab.NewClient(nil, cfg.GitlabToken)
	err := git.SetBaseURL(cfg.GitlabBaseURL)
	if err != nil {
		log.Fatalf("Can not set gitlab base url: %v", err)
	}

	db, err := sql.Open("sqlite3", cfg.DbPath)
	if err != nil {
		log.Fatalf("Can not open sqlite3 database with path %q: %v", cfg.DbPath, err)
	}
	defer db.Close()

	// mengine := newMigrationEngine(db)
	// err = mengine.run()
	// if err != nil {
	// 	log.Fatalf("Can not apply migrations: %v", err)
	// }

	/*
		bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
		if err != nil {
			log.Fatalf("Can not initialize telegram bot: %v", err)
		}
		bot.Debug = true
		log.Printf("Authorized on account %s", bot.Self.UserName)
	*/
	brain := review.NewBrain(&dummylister.Lister{})
	reviewHandlers := &web.ReviewHandlers{Brain: brain}

	mux := http.NewServeMux()
	server := web.NewServer(mux, reviewHandlers)
	server.RegisterRoutes()

	/*
		webServer := web.NewServer(bot)
		webServer.RegisterRoutes(mux)
		webServer.RunTelegram()
	*/

	log.Println("Listen http://0.0.0.0:6789")
	err = server.Listen("0.0.0.0:6789")
	if err != nil {
		log.Fatalf("Web server error: %v", err)
	}
}
