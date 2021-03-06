package cmd

import (
	"fmt"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type LinkedID struct {
	UserID  int
	GroupID int64
}

func Bot(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	weatherChan := make(chan string)

	updates, err := bot.GetUpdatesChan(u)

	unames := map[LinkedID]bool{}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() != true {
			if _, ok := unames[LinkedID{UserID: update.Message.From.ID, GroupID: update.Message.Chat.ID}]; ok {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				go getWeather(update.Message.Chat.FirstName, update.Message.Chat.LastName, update.Message.Text, weatherChan)
				msg.Text = <-weatherChan
				bot.Send(msg)
				delete(unames, LinkedID{update.Message.From.ID, update.Message.Chat.ID})
			}
		} else if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			delete(unames, LinkedID{UserID: update.Message.From.ID, GroupID: update.Message.Chat.ID})
			switch update.Message.Command() {
			case "start":
				msg.Text = "Hello! Welcome to Astronomia Bot. Type /help to see all available commands."
			case "help":
				msg.Text = "Type /sayhi to say hi.\nType /status to get bot status.\nType /weather to get current weather information.\nType /help to see all available commands."
			case "sayhi":
				if update.Message.Chat.IsGroup() {
					msg.Text = fmt.Sprintf("Hi %s!", update.Message.Chat.UserName)
				} else {
					msg.Text = fmt.Sprintf("Hi %s %s! Have a good day!", update.Message.Chat.FirstName, update.Message.Chat.LastName)
				}
			case "status":
				msg.Text = "I'm ok."
			case "weather":
				if _, ok := unames[LinkedID{update.Message.From.ID, update.Message.Chat.ID}]; !ok {
					unames[LinkedID{update.Message.From.ID, update.Message.Chat.ID}] = true
				}
				msg.Text = "What is the address you want to know the weather of?"
			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)
		}
	}
}

func getWeather(firstName, lastName, address string, weatherChan chan string) {
	msg := GetWeather(firstName, lastName, address)
	weatherChan <- msg
}
