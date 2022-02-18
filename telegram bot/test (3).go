package main

import (
	"fmt"
	"time"

	//"io/ioutil"

	//"v/ntlmssp"

	"github.com/dgrijalva/jwt-go"
)

type SupportClient struct {
}

type Sth struct {
	Value interface{}
}

type Resp struct {
	ProcessDefinitionKey, State, Id string
}

type Sth2 struct {
	Value string
}

var m map[string]interface{}

func main() {
	// c := &http.Client{
	// 	Transport: ntlmssp.Negotiator{
	// 		RoundTripper: &http.Transport{},
	// 	},
	// }

	// var dst bytes.Buffer
	// url := "http://halykbpm-auth.halykbank.nb/WindowsAuthentication/auth/bearer?clientId=spmapi"
	// req, err := http.NewRequest(http.MethodGet, url, nil)
	// req.SetBasicAuth("00052920", "Xanx@123")
	// res, err := c.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// buffer := make([]byte, 1)
	// fmt.Println("Print body before closing")
	// bytess, err := io.CopyBuffer(&dst, res.Body, buffer)
	// if err != nil {

	// 	// if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
	// 	log.Fatal(err)
	// }
	// res.Body.Close()
	// fmt.Println(bytess)
	// fmt.Println("string of dst bytes", string(dst.Bytes()))
	// var token string
	// if err := json.Unmarshal(dst.Bytes(), &token); err != nil {
	// 	fmt.Println("unmarshal error", err)
	// 	return
	// }

	// fmt.Println("token\n", token)
	// // 2nd request

	// // reqz.Header.Set("Authorization", "NTLM TlRMTVNTUAABAAAAB6IIogAAAAAoAAAAAAAAACgAAAAFASgKAAAADw==")

	// // bod, _ := ioutil.ReadAll(res.Body)
	// // fmt.Println("body", string(bod))
	// // // fmt.Println(string(bod), "Header:\n", respv.Header)
	// // token := string(bod)

	// IIN := "841227351343" // ge09bc30f1506015; 33 {android}; 37 done_success 125 {https://fr.halykbank.kz}  126 {d06284e0-6019-4812-b9a8-c2ecb7eecd09}  127 {done_success}
	// //IIN = "891017301067" // ID: gfb0401718742022; 34 {android}; 38 done_canceled; 132 {0d8e11ab-b665-40f1-874f-289e8eef3c69}; 133 {done_cancelled}
	// IIN = "860202350796" // gd5f31eeed713014; 33-android 37, 127 - done_success 125 {https://fr.halykbank.kz}126 {1903151c-6449-4c2a-af51-391e073d16f5}
	// IIN = "970906000196" // gfafb66b84726024; 34 {iOS}; 38 unknown_status; max 93 values!!! no values in range from 127
	// // c := &http.Client{}
	// getProcessesURL := "https://halykbpm-api.halykbank.nb/process-searcher/instance?searchValue=" + IIN

	// //token := "eyJhbGciOiJodHRwOi8vd3d3LnczLm9yZy8yMDAxLzA0L3htbGRzaWctbW9yZSNobWFjLXNoYTI1NiIsInR5cCI6IkpXVCJ9.eyJodHRwOi8vc2NoZW1hcy54bWxzb2FwLm9yZy93cy8yMDA1LzA1L2lkZW50aXR5L2NsYWltcy9uYW1lIjoiMDAwNTI5MjAiLCJhZGxvZ2luIjoiVU5JVkVSU0FMXFwwMDA1MjkyMCIsImV4cCI6MTY0MzgyMzQ5NywiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDoxNTM1L2F1dGgvaXNzdWVyIiwiYXVkIjoic3BtYXBpIn0.37Rl3YJ_POCCwh6VxgAhRxd6VW9kdqT8-9uzm1OsTJc"
	// req, err = http.NewRequest("GET", getProcessesURL, nil)
	// req.Header.Set("Authorization", "Bearer "+token)
	// res, err = c.Do(req)
	// if err != nil {
	// 	panic(err)
	// }

	// //    bod, err := ioutil.ReadAll(res.Body)
	// //    if err != nil {
	// //       fmt.Println("bod err", err)
	// //       return
	// //    }
	// // fmt.Println("bod", string(bod), "\n\n\n")

	// dst.Reset()
	// bytess, err = io.CopyBuffer(&dst, res.Body, buffer)
	// if err != nil {

	// 	// if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
	// 	log.Fatal(err)
	// }
	// res.Body.Close()
	// fmt.Println(bytess)

	// //body, _ := ioutil.ReadAll(res.Body)
	// // var m map[string]interface{}
	// // if err := json.Unmarshal(dst.Bytes(), &m); err != nil {
	// //    fmt.Println("unmarshal error", err)
	// //    return
	// // }
	// // for key, v := range m {
	// //    fmt.Println(key,v,"key,value")
	// // }
	// processNeeded := "onboarding01"

	// fmt.Println("dst bytes:\n\n\n\n\n", string(dst.Bytes()))
	// var B []Resp
	// if err := json.Unmarshal(dst.Bytes(), &B); err != nil {
	// 	fmt.Println("Unramshal json error", err)
	// 	return
	// }
	// var processID string
	// // TODO: separate function for this: else return false!
	// for _, process := range B {
	// 	if process.ProcessDefinitionKey == processNeeded {
	// 		processID = process.Id
	// 		fmt.Println("Printing status", process.State)
	// 		if process.State == "ACTIVE" {
	// 			fmt.Println("Process is still active")
	// 			return
	// 		}
	// 		break
	// 	}
	// }

	// fmt.Println(processID)
	// fmt.Println("response sts", res.Status)
	// //processID := B[0].Id //-- любой посл процесс

	// getProcessVarURL := "https://halykbpm-api.halykbank.nb/bpm-front-webapi/api/history/variable-instance?"

	// // //var jsonStr = []byte(`{"processInstanceIdIn":["gd8894a224281016"]}`)
	// dst.Reset()
	// var jsonStr = []byte(fmt.Sprintf("{\"processInstanceIdIn\":[\"%s\"]}", processID))
	// // var jsonStr = []byte(fmt.Sprintf("{\"processInstanceIdIn\":[\"%s\"]}", processID))
	// req, err = http.NewRequest("POST", getProcessVarURL, bytes.NewBuffer(jsonStr))
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer "+token)

	// //client := &http.Client{}
	// //resp2, err := client.Do(req2)
	// res, err = c.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// //defer res.Body.Close()
	// //body, _ = ioutil.ReadAll(res.Body)
	// bytess, err = io.CopyBuffer(&dst, res.Body, buffer)
	// if err != nil {

	// 	// if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
	// 	log.Fatal(err)
	// }
	// res.Body.Close()
	// fmt.Println(bytess)
	// var A []Sth
	// //var A []Sth2
	// if err := json.Unmarshal(dst.Bytes(), &A); err != nil {
	// 	fmt.Println("err last req unmarshal", err)
	// 	return
	// }
	// for i, hing := range A {
	// 	fmt.Println(i, hing)
	// }
	// check := false
	// for i := 35; i <= 40; i++ {
	// 	message := A[i]
	// 	if value, ok := message.Value.(string); ok && len(value) > 4 {
	// 		//if value := message.Value; len(value)>4 && value[:4] == "done" {
	// 		fmt.Println(message)
	// 		check = true
	// 		break
	// 	}
	// }
	// if check == false {
	// 	fmt.Println("No process was found")
	// }
	// //fmt.Println("response Body:", string(body))
	// fmt.Println("response sts", res.Status)
	// // bytes, err = io.CopyBuffer(dst, res.Body, buffer)
	// // if err != nil {

	// // 	// if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
	// // 	log.Fatal(err)
	// // }

	// fmt.Println(bytess)
	var tokenString = "eyJhbGciOiJodHRwOi8vd3d3LnczLm9yZy8yMDAxLzA0L3htbGRzaWctbW9yZSNobWFjLXNoYTI1NiIsInR5cCI6IkpXVCJ9.eyJodHRwOi8vc2NoZW1hcy54bWxzb2FwLm9yZy93cy8yMDA1LzA1L2lkZW50aXR5L2NsYWltcy9uYW1lIjoiMDAwNTI5MjAiLCJhZGxvZ2luIjoiVU5JVkVSU0FMXFwwMDA1MjkyMCIsImV4cCI6MTY0NDgzODEwNSwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDoxNTM1L2F1dGgvaXNzdWVyIiwiYXVkIjoic3BtYXBpIn0.prazUWjhj5EkJUYFAFMdVtYMXgmJyEtCtuG7x-bg-vc"

	tokenn, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err, ok := err.(*jwt.ValidationError); ok && err.Errors != jwt.ValidationErrorUnverifiable || !ok {

		//} != "signing method (alg) is unavailable." {
		fmt.Println(err)
		return
	}

	if claims, ok := tokenn.Claims.(jwt.MapClaims); ok {
		exp, ok := claims["exp"].(float64)
		if !ok {
			fmt.Println("No")
			return
		}
		expiredTime := time.Unix(int64(exp), 0)
		if time.Now().After(expiredTime) {
			fmt.Println("expired")
			return
		}

	} else {
		fmt.Println(err)
	}
}

// func IOCopy(reader io.Reader) (map[string]interface{}, error) {
//    var (
//       m    map[string]interface{}
//       buf  bytes.Buffer
//       _, _ = io.Copy(&buf, reader)
//    )

//    return m, json.Unmarshal(buf.Bytes(), &m)
// }

//    resp, err := http.Get("halykbpm-auth.halykbank.nb/WindowsAuthentication/auth/bearer?clientId=spmapi")
//    if err != nil {
//       log.Fatalln(err)
//    }
// //We Read the response body on the line below.
//    body, err := ioutil.ReadAll(resp.Body)
//    if err != nil {
//       log.Fatalln(err)
//    }
// //Convert the body to type string
//    sb := string(body)
//    log.Printf(sb)
// }

// package main

// import (
//    "io/ioutil"
//    "fmt"
//    "net/http"
//    "bytes"
//    "encoding/json"
//    "v/ntlmssp"
//    //"github.com/Azure/go-ntlmssp"
// )

// type Sth struct {
//    Value interface{}
// }

// type Resp struct {
//    ProcessDefinitionKey, Id string
// }

// func main() {
//       c := &http.Client{
//       Transport: ntlmssp.Negotiator{
//           RoundTripper: &http.Transport{},
//       },
//   }
//    url := "http://halykbpm-auth.halykbank.nb/WindowsAuthentication/auth/bearer?clientId=spmapi"
//    reqz, err := http.NewRequest("GET",  url, nil)

//   // reqz.Header.Set("Authorization", "NTLM TlRMTVNTUAABAAAAB6IIogAAAAAoAAAAAAAAACgAAAAFASgKAAAADw==")
//    reqz.SetBasicAuth("00052920", "Xanx@123")

//    //respv, err := cc.Do(reqz)
//    respv, err := c.Do(reqz)
//    if err != nil {
//       panic(err)
//    }

//    defer respv.Body.Close()
//    bod, _ := ioutil.ReadAll(respv.Body)
//    fmt.Println("body", string(bod))
//   // fmt.Println(string(bod), "Header:\n", respv.Header)
//    token := string(bod)
//    token = token[1:len(bod)-1]
//    fmt.Println("token\n",token)
//     IIN := "841227351343"
//    // c := &http.Client{}
//     getProcessesURL := "https://halykbpm-api.halykbank.nb/process-searcher/instance?searchValue=" + IIN

//      //token := "eyJhbGciOiJodHRwOi8vd3d3LnczLm9yZy8yMDAxLzA0L3htbGRzaWctbW9yZSNobWFjLXNoYTI1NiIsInR5cCI6IkpXVCJ9.eyJodHRwOi8vc2NoZW1hcy54bWxzb2FwLm9yZy93cy8yMDA1LzA1L2lkZW50aXR5L2NsYWltcy9uYW1lIjoiMDAwNTI5MjAiLCJhZGxvZ2luIjoiVU5JVkVSU0FMXFwwMDA1MjkyMCIsImV4cCI6MTY0MzgyMzQ5NywiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDoxNTM1L2F1dGgvaXNzdWVyIiwiYXVkIjoic3BtYXBpIn0.37Rl3YJ_POCCwh6VxgAhRxd6VW9kdqT8-9uzm1OsTJc"
//      req, err := http.NewRequest("GET", getProcessesURL, nil)
//      req.Header.Set("Authorization", "Bearer "+token)
//      resp, err := c.Do(req)
//      if err != nil {
//          panic(err)
//      }

//      body, _ := ioutil.ReadAll(resp.Body)
//      var B []Resp
//      if err := json.Unmarshal(body, &B); err != nil {
//      	fmt.Println(err)
//      	return
//      }
//      fmt.Println(B[0].Id)
//      fmt.Println("response sts", resp.Status)
//    processID := B[0].Id
//     defer resp.Body.Close()
//      getProcessVarURL := "https://halykbpm-api.halykbank.nb/bpm-front-webapi/api/history/variable-instance?"

//      // //var jsonStr = []byte(`{"processInstanceIdIn":["gd8894a224281016"]}`)
//      var jsonStr = []byte(fmt.Sprintf("{\"processInstanceIdIn\":[\"%s\"]}", processID))
//     //var jsonStr = []byte(`{"processInstanceIdIn":["gf279d39f789c024"]}`)
//      req2, err := http.NewRequest("POST", getProcessVarURL, bytes.NewBuffer(jsonStr))
//      req2.Header.Set("Content-Type", "application/json")
//      req2.Header.Set("Authorization", "Bearer "+token)

//      //client := &http.Client{}
//      //resp2, err := client.Do(req2)
//      resp2, err := c.Do(req2)
//      if err != nil {
//          panic(err)
//      }
//      defer resp2.Body.Close()
//      body, _ = ioutil.ReadAll(resp2.Body)
//      var A []Sth
//      if err := json.Unmarshal(body, &A); err != nil {
//         fmt.Println(err)
//         return
//      }
//      fmt.Println(A[37])
//      //fmt.Println("response Body:", string(body))
//      fmt.Println("response sts", resp2.Status)
// }

//    resp, err := http.Get("halykbpm-auth.halykbank.nb/WindowsAuthentication/auth/bearer?clientId=spmapi")
//    if err != nil {
//       log.Fatalln(err)
//    }
// //We Read the response body on the line below.
//    body, err := ioutil.ReadAll(resp.Body)
//    if err != nil {
//       log.Fatalln(err)
//    }
// //Convert the body to type string
//    sb := string(body)
//    log.Printf(sb)
// }

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strings"
// )

// func main() {
// 	client := &http.Client{}
// 	body := strings.NewReader("{\"processInstanceIdIn\": [\"gd8894a224281016\"]")
// 	req, err := http.NewRequest("POST", "https://halykbpm-api.halykbank.nb/bpm-front-webapi/api/history/variable-instance?", body)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	req.Header.Add("Content-Type", "application/json")
//    token := "eyJhbGciOiJodHRwOi8vd3d3LnczLm9yZy8yMDAxLzA0L3htbGRzaWctbW9yZSNobWFjLXNoYTI1NiIsInR5cCI6IkpXVCJ9.eyJodHRwOi8vc2NoZW1hcy54bWxzb2FwLm9yZy93cy8yMDA1LzA1L2lkZW50aXR5L2NsYWltcy9uYW1lIjoiMDAwNTI5MjAiLCJhZGxvZ2luIjoiVU5JVkVSU0FMXFwwMDA1MjkyMCIsImV4cCI6MTY0MzgxNzA3NCwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDoxNTM1L2F1dGgvaXNzdWVyIiwiYXVkIjoic3BtYXBpIn0.GyrGHa7ejwwDjUtGtA9pJuDR-tnOcxN8vdX9ydLuZvI"
// 	req.Header.Add("Authorization", "Bearer "+token)
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	fmt.Println(resp.Status)
// }
