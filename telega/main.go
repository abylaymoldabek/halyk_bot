// package main

// import (
// 	"log"
// 	// "os"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// func main() {
// 	bot, err := tgbotapi.NewBotAPI("5001533822:AAHqehWoBVXpqiSwXMq3i9GX4znSw0D3d9s")
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	bot.Debug = true

// 	log.Printf("Authorized on account %s", bot.Self.UserName)

// 	u := tgbotapi.NewUpdate(0)
// 	u.Timeout = 60

// 	updates := bot.GetUpdatesChan(u)

// 	for update := range updates {
// 		if update.Message == nil { // ignore any non-Message updates
// 			continue
// 		}

// 		if !update.Message.IsCommand() { // ignore any non-command Messages
// 			continue
// 		}

// 		// Create a new MessageConfig. We don't have text yet,
// 		// so we leave it empty.
// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

// 		// Extract the command from the Message.
// 		switch update.Message.Command() {
// 		case "help":
// 			msg.Text = "I understand /sayhi and /status."
// 		case "sayhi":
// 			msg.Text = "Hi :)"
// 		case "status":
// 			msg.Text = "I'm ok."
// 		default:
// 			msg.Text = "I don't know that command"
// 		}

// 		if _, err := bot.Send(msg); err != nil {
// 			log.Panic(err)
// 		}
// 	}
// }

package main

import (
	"log"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strings"
)

func main() {
	sampleRegexp := regexp.MustCompile(`\d`)
	bot, err := tgbotapi.NewBotAPI("5001533822:AAHqehWoBVXpqiSwXMq3i9GX4znSw0D3d9s")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	teMessage := ""
	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			if strings.Contains(update.Message.Text, "/") {
				continue
			} else if sampleRegexp.MatchString(update.Message.Text) {
				teMessage += update.Message.Text
				fmt.Println(teMessage)
				msg.Text = "Хорошо, получил данные. Прошу ожидайте..."
				bot.Send(msg)
			} else {
				msg.Text = "Неправильные данные"	
				bot.Send(msg)
			}
		}
	}
}
