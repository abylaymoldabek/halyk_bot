package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"v/domain"
	"v/process/repository/client"
	"v/process/usecase"
	"time"
	"sync"
	"log"
	"regexp"
	tgbotapi "v/telegram"
	
	"strings"

	//"github.com/subosito/gotenv"
)

// func init() {
// 	gotenv.Load()
// }

func SetEnvAll() {
	os.Setenv("CTX_TIMEOUT", "500000000")
	os.Setenv("TOKEN_URL","http://halykbpm-auth.halykbank.nb/WindowsAuthentication/auth/bearer?clientId=spmapi")
	os.Setenv("USERNAME","00052920")
	os.Setenv("PASSWORD","Xanx@123")
	os.Setenv("PROCESSES_URL","https://halykbpm-api.halykbank.nb/process-searcher/instance?searchValue=")
	os.Setenv("PROCESS_URL","https://halykbpm-api.halykbank.nb/bpm-front-webapi/api/history/variable-instance?")
	os.Setenv("GET_INCIDENT_URL","https://halykbpm-api.halykbank.nb/bpm-front-webapi/api/incident?processInstanceId=")
	os.Setenv("RETRY_JOB_URL","https://halykbpm-api.halykbank.nb/bpm-front-webapi/api/job")
	os.Setenv("RETRY_TASK_URL","https://halykbpm-api.halykbank.nb/bpm-front-webapi/api/external-task/retries")



}

func main() {
	SetEnvAll() 
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
	
	processRepo := client.NewClient()
	time.Sleep(time.Second*20)
	timeoutInt, err := strconv.Atoi(os.Getenv("CTX_TIMEOUT"))
	if err != nil {
		fmt.Println("Invalid timeout")
		return
	}
	timeoutContext := time.Duration(timeoutInt) * time.Second
	pu := usecase.NewProcessUsecase(processRepo, timeoutContext)
	criteria := domain.Criteria{
		ID: teMessage,   //"790713303493",ration is nill "credit01"
		Type: "onboarding01",
	}

	
	res, err := pu.MainLogic(context.Background(), criteria)
	if err != nil {
		fmt.Println(err)
		return
	}
}
