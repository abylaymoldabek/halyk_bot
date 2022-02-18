package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"support/domain"
	"support/process/repository/api"
	"support/process/usecase"
	"time"

	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

func main() {
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
		ID:   "",
		Type: "onboarding01",
	}
	pu.MainLogic(context.Background(), criteria)
}
