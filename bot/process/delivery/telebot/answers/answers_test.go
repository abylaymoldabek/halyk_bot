package answers_test

import (
	"log"
	"testing"
	"v/process/delivery/telebot/answers"
)

func TestAnswers_NewAnswers(t *testing.T) {
	testCases := []struct {
		status  string
		isValid bool
	}{
		{
			status: "dm1_decline",
			isValid: true,
		},
		{
			status: "done_rejected_publicId",
			isValid: true,
		},
		{
			status: "qwerty",
			isValid: false,
		},
	}
	ans := answers.NewAnswers()
	
	for _, tc := range testCases {
		text, ok := ans[tc.status]
		if tc.isValid != ok {
			t.Error("error")
		}
		log.Println(text)
	}
}

// func TestAnswers_GetAnswer(t *testing.T) {
// 	ans := answers.NewAnswers()

// 	res, err := ans.GetAnswer("done_success")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	log.Println(res)
// }
