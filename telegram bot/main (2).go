package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"v/domain"
	"v/process/repository/api"
	"v/process/usecase"
	"time"
	"sync"

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
	processRepo := api.NewClient()
	timeoutInt, err := strconv.Atoi(os.Getenv("CTX_TIMEOUT"))
	if err != nil {
		fmt.Println("Invalid timeout")
		return
	}
	timeoutContext := time.Duration(timeoutInt) * time.Second
	pu := usecase.NewProcessUsecase(processRepo, timeoutContext)
	// _articleHttpDelivery.NewArticleHandler(e, au)
	criteria := domain.Criteria{
		ID: "gf9645b3b7233023",   //"790713303493",ration is nill "credit01"
		Type: "unitedCredit",
	}

	var wg sync.WaitGroup



  
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int, c *api.Client) {
			defer wg.Done()
			fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
			res, err := pu.MainLogic(context.Background(), criteria)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(i, res)
			//fmt.Println("tokenSTRING", i, c.Token.tokenString)
		}(i, processRepo)
	}
	wg.Wait()
}
