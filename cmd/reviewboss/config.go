package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	GitlabToken      string
	GitlabBaseURL    string
	DbPath           string
	TelegramBotToken string
}

func loadConfig() *config {
	cfg := &config{}

	// set token
	cfg.GitlabToken = os.Getenv("GITLAB_TOKEN")

	// set baseurl
	url, ok := os.LookupEnv("GITLAB_BASEURL")
	if !ok {
		url = "https://gitlab.com/api/v4"
	}
	cfg.GitlabBaseURL = url

	// set path to db
	path, ok := os.LookupEnv("GITLAB_BASEURL")
	if !ok {
		path = "data.db"
	}
	cfg.DbPath = path

	// telegram
	cfg.TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")

	return cfg
}
