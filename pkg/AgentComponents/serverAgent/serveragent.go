package serverAgent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Task struct {
	Id             int    `json:"id"`
	Arg1           string `json:"arg1"`
	Arg2           string `json:"arg2"`
	Operation      string `json:"operation"`
	Operation_time int    `json:"operation_time"`
}
type Result struct {
	Id     int    `json:"id"`
	Result string `json:"result"`
}

func AgentStart() {
	go getRequest()
	workerPool()
}

var text Task
var ch = make(chan Task, 5)

func getRequest() {
	url := "http://localhost:8080/internal/task/"
	for {
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&text); err != nil {
			panic(err)
		}
		if text.Id == 0 {
			continue
		}
		ch <- text
		fmt.Println(text)
	}
}
func workerPool() {
	for {
		miniTask := <-ch
		time.Sleep(time.Duration(miniTask.Operation_time * int(time.Millisecond)))
		fmt.Println("Отправляем на расчет", miniTask)
		go ParalelCalc(miniTask) // Результат вычисления
	}
}
func PostResult(res Result) {
	url := "http://localhost:8080/internal/task/"
	jsonData, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		panic(err)
	}
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	r.Header.Set("Content-Type", "application/json")
	fmt.Println("Post Запрос ", res)
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
}
