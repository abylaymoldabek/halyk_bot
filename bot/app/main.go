package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"v/logger"
	"v/process/delivery/telebot"
	"v/process/repository/client"
	"v/process/usecase"
	//"github.com/subosito/gotenv"
)

// func init() {
// 	gotenv.Load()
// }

func SetEnvAll() {
	os.Setenv("CTX_TIMEOUT", "500000000")
	os.Setenv("TOKEN_URL", "http://halykbpm-auth.halykbank.nb/WindowsAuthentication/auth/bearer?clientId=spmapi")
	os.Setenv("USERNAME", "00052920")
	os.Setenv("PASSWORD", "Xanx@123")
	os.Setenv("PROCESSES_URL", "https://halykbpm-api.halykbank.nb/process-searcher/instance?searchValue=")
	os.Setenv("PROCESS_URL", "https://halykbpm-api.halykbank.nb/bpm-front-webapi/api/history/variable-instance?")
	os.Setenv("GET_INCIDENT_URL", "https://halykbpm-api.halykbank.nb/bpm-front-webapi/api/incident?processInstanceId=")
	os.Setenv("RETRY_JOB_URL", "https://halykbpm-api.halykbank.nb/bpm-front-webapi/api/job")
	os.Setenv("RETRY_TASK_URL", "https://halykbpm-api.halykbank.nb/bpm-front-webapi/api/external-task/retries")

}

func main() {
	SetEnvAll()
	processRepo := client.NewClient()
	//time.Sleep(time.Second * 20)
	timeoutInt, err := strconv.Atoi(os.Getenv("CTX_TIMEOUT"))
	if err != nil {
		fmt.Println("Invalid timeout")
		return
	}
	log := logger.NewLogger()
	timeout := time.Duration(timeoutInt) * time.Second
	process_uc := usecase.NewProcessUsecase(processRepo, timeout)
	process_handler := telebot.NewProcessHandler(log, process_uc)
	process_handler.ProcessRequest()
}
