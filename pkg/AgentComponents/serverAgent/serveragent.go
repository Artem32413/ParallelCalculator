package serveragent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	c "parallelcalculator/pkg/AgentComponents/calculation"
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
var ch = make(chan Task, 2)
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
		// fmt.Println("Агент добавляет в очередь ", text)
		ch <- text
		fmt.Println(text)
	}
}
func workerPool() {
	for {
		miniTask := <-ch
		// fmt.Println("Агент начинает вычисление ", miniTask)
		time.Sleep(time.Duration(miniTask.Operation_time * int(time.Second)))
		eq := c.Calcularion(miniTask.Arg1, miniTask.Operation, miniTask.Arg2) // Результат вычисления
		// fmt.Println("Ответ = ", eq)
		var res = Result{
			Id:     miniTask.Id,
			Result: eq,
		}
		postResult(res)
	}
}
func postResult(res Result) {
    url := "http://localhost:8080/internal/task/"
    jsonData, err := json.MarshalIndent(res, "", "    ")
    if err != nil {
        panic(err)
    }
    // fmt.Println("Отправляем данные:", string(jsonData))
    
    r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        panic(err)
    }
    r.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    response, err := client.Do(r)
    if err != nil {
        panic(err)
    }
    defer response.Body.Close()
    
    // fmt.Println("Ответ от сервера:", response.Status)
}
