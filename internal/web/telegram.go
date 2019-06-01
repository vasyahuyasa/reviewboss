package web

import (
	"log"
	"net/http"
)

func telegramHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request from telegram:", r)
}

func telegramUpdate(telegram.)