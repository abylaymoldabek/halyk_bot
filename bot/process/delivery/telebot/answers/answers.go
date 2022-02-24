package answers

import (
	"encoding/json"
	_ "errors"
	"io/ioutil"
)

const (
	filePath = "./answers.json"
)

type Answers map[string]string

func NewAnswers() Answers {
	var answers Answers

	file, _ := ioutil.ReadFile(filePath)
	if err := json.Unmarshal(file, &answers); err != nil {
		panic(err)
	}

	return answers
}

// type Case struct {
// 	Status string
// 	Answer string
// }

// type ProcessInfo struct {
// 	ProcessName string
// 	Cases       []Case
// }

// type Answers []ProcessInfo


// func (answers Answers) GetAnswer(processName, processStatus string) (string, error) {
// 	for _, ans := range answers {
// 		if ans.ProcessName == processName {
// 			for _, cs := range ans.Cases {
// 				if cs.Status == processStatus {
// 					return cs.Answer, nil
// 				}
// 			}
// 		}
// 	}

// 	return "", errors.New("НЕТ ТАКОГО КЕЙСА")
// }