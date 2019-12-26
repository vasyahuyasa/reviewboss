package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/proxy"
)

type botAPI struct {
	token     string
	proxyAddr string
	tgbot     *tgbotapi.BotAPI
	onUpdate  func(tgbotapi.Update)
}

func (api *botAPI) waitForUpdates() error {
	var httpClient *http.Client

	if api.proxyAddr != "" {
		dialer, err := proxy.SOCKS5("tcp", api.proxyAddr, nil, nil)
		if err != nil {
			return fmt.Errorf("can't connect to the proxy: %w", err)
		}

		httpClient = &http.Client{Transport: &http.Transport{Dial: dialer.Dial}}
	} else {
		httpClient = &http.Client{}
	}

	bot, err := tgbotapi.NewBotAPIWithClient(api.token, httpClient)
	if err != nil {
		return fmt.Errorf("can not inittialize telegram bot: %w", err)
	}

	log.Printf("telegram bot authorized on account %s", bot.Self.UserName)

	api.tgbot = bot

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := api.tgbot.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("telegram can not get update channel: %w", err)
	}

	// clear old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	go func() {
		for update := range updates {
			api.onUpdate(update)
		}
	}()

	return nil
}
