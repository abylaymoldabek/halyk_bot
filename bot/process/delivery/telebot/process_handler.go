package telebot

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"v/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ProcessHandler struct {
	log domain.Logger
	uc  domain.ProcessUsecase
	reg *regexp.Regexp
}

func NewProcessHandler(log domain.Logger, uc domain.ProcessUsecase) *ProcessHandler {
	return &ProcessHandler{
		log: log,
		uc:  uc,
		reg: regexp.MustCompile(`\d`),
	}
}

func (p *ProcessHandler) ProcessRequest() {

	bot, err := tgbotapi.NewBotAPI("5001533822:AAHqehWoBVXpqiSwXMq3i9GX4znSw0D3d9s") // better export
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	p.log.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	teMessage := ""
	for update := range updates {
		if update.Message != nil { // If we got a message
			p.log.Info(fmt.Sprintf("[%s] %s", update.Message.From.UserName, update.Message.Text))
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			if strings.Contains(update.Message.Text, "/") {
				continue
			} else if p.reg.MatchString(update.Message.Text) {
				teMessage += update.Message.Text
				fmt.Println(teMessage)
				msg.Text = "Хорошо, получил данные. Прошу, ожидайте..."
				bot.Send(msg)

				criteria := domain.Criteria{
					ID:   teMessage,      // иин/номер телефона/айди процесса
					Type: "onboarding01", // replace with variable storing category
				}
				// добавить еще перевод:
				// если категория перевод то дальнейшие вопросы (Диас составит часто задаваемые и на них стандартный ответ)
				// эти вопросы тоже как боксы сделать чтоб на них нажимали и ответы появлялись
				// если какой то уникальный вопрос то сказать обращайтесь в поддержку
				// если просят что то поменять (тоже в боксы добавить: initRole поменять, направить на УВК)
				// то вызываем другую функцию которую я седня напишу: e.g.
				// if request == "UVK" {
				// 	res, err := p.uc.ProcessTransfer(criteria)
				// }
				// в criteria добавим еще филиал

				res, err := p.uc.ProcessRequest(context.Background(), p.log, criteria)
				if err != nil {
					p.log.Error(err, "сюда табельный")
					return
				}
				msg.Text = res.Status
				bot.Send(msg)
			} else {
				msg.Text = "Неправильные данные"
				bot.Send(msg)
			}
		}
	}
}
