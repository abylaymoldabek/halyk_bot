package answers_test

import (
	"log"
	"testing"
	"v/process/delivery/telebot/answers"
)

func TestAnswers_NewAnswers(t *testing.T) {
	ans := answers.NewAnswers()
	log.Println(len(ans))
	log.Println(ans[0].ProcessName)
	log.Println(ans[0].Cases)
}

func TestAnswers_GetAnswer(t *testing.T) {
	ans := answers.NewAnswers()

	res, err := ans.GetAnswer("ONB", "done_success")
	if err != nil {
		t.Error(err)
	}
	log.Println(res)
}