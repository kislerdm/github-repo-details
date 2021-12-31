package main

import (
	"log"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *telegram.BotAPI

func init() {
	var err error
	bot, err = telegram.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}
	bot.Debug = true
}

func main() {
	u := telegram.NewUpdate(0)
	u.Timeout = 10

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}

		msg := telegram.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand and /start."
		case "start":
			msg.Text = `<pre>|a|b|
|:-:|:-:|
|[a](https://google.com)|1|
</pre>`
		}

		if _, err := bot.Send(msg); err != nil {
			log.Fatalln(err)
		}
	}
}
