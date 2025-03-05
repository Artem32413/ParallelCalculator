package serverorcestrator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	mu         sync.Mutex
	Ids        IdExpressions
	expression Expressions
	e          error
)
var math_expr = make(map[int]MathExpr) // Все принимаемые выражения
var m = make(map[int]Expressions)      // Результаты вычисления
var ch = make(chan int, 5)             // Канал для воркер пула

type IdExpressions struct { // Глобальный счетчик Id
	Id int `json:"id"`
}
type MathExpr struct { // принимаемое выражение
	Expression string `json:"expression"`
}

type Expressions struct { // Результат выражения от Агента
	Id     int    `json:"id"`
	Status string `json:"status"`
	Result string `json:"result"`
}
type Task struct { // Разбивка на действия Агенту
	Id             int    `json:"id"`
	Arg1           string `json:"arg1"`
	Arg2           string `json:"arg2"`
	Operation      string `json:"operation"`
	Operation_time int    `json:"operation_time"`
}

func StartServer() {
	go workerPool()
	http.HandleFunc("/api/v1/calculate/", calculate)
	http.HandleFunc("/api/v1/expressions/", expressions)
	http.HandleFunc("/api/v1/expressions/{id}", expressionsId)
	http.HandleFunc("/internal/task/", task)
	http.ListenAndServe(":8080", nil)
}
func calculate(w http.ResponseWriter, r *http.Request) {
	// Примем выражение для расчета, вернем id и сохраним в мапу => m
	var mathEx MathExpr
	// Читаем пришедшее выражение
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&mathEx); err != nil {
		http.Error(w, "Ошибка C декодингом", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err := json.MarshalIndent(mathEx, "", "    ")
	if err != nil {
		http.Error(w, "Ошибка при преобразовании в JSON:", http.StatusInternalServerError)
		return
	}
	expression := string(mathEx.Expression)
	expression = strings.ReplaceAll(expression, " ", "")
	if err = ValidateExpression(expression); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	mu.Lock()
	Ids.Id++
	m[Ids.Id] = Expressions{
		Id:     Ids.Id,
		Status: "принято",
		Result: "",
	}
	math_expr[Ids.Id] = MathExpr{
		Expression: mathEx.Expression,
	}
	mu.Unlock()
	// Отправим ответ id
	var idExpr IdExpressions
	idExpr.Id = Ids.Id
	jsonDataId, err := json.MarshalIndent(idExpr, "", "    ")
	if err != nil {
		http.Error(w, "Ошибка при маршелинге:", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonDataId)
	ch <- Ids.Id
	// Запускаем для расчета
}

func workerPool() {
	for {
		idTask := <-ch
		go Priority(idTask, []byte(math_expr[idTask].Expression))
	}
}
func task(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet { // Вернем задание агенту
		ti := GetTask()
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.MarshalIndent(ti, "", "    ")
		if err != nil {
			http.Error(w, "Ошибка при преобразовании в JSON:", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)

	} else if r.Method == http.MethodPost { // Принимаем ответ от агента
		var t2 R
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&t2); err != nil {
			http.Error(w, "Ошибка с декодингом", http.StatusInternalServerError)
			return
		}
		if t2.Id == 0 {
			http.Error(w, "Некорректный ID задачи", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Println("Начало блока Task")
		mux.Lock()
		RPerem[t2.Id] = t2
		mux.Unlock()
		fmt.Println("Конец блока Task")
	}
	defer r.Body.Close()
}

func GetTask() Task { // Вернет одно действие на расчет
	var ti Task
	for _, el := range Rtmp {
		ti = Task{
			Id:             el.Id,
			Arg1:           el.Num1,
			Arg2:           el.Num2,
			Operation:      el.Operator,
			Operation_time: TimeSleep(el.Operator),
		}

		delete(Rtmp, el.Id) // Удаляем из map

		return ti // Возвращаем рандомное действие
	}
	return ti // Заглушка
}
func expressions(w http.ResponseWriter, r *http.Request) {
	var exp []Expressions
	mux.Lock()
	for _, el := range m {
		exp = append(exp, el)
	}
	mux.Unlock()
	sort.Slice(exp, func(i, j int) bool {
		return exp[i].Id < exp[j].Id
	})
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.MarshalIndent(exp, "", "    ")
	if err != nil {
		http.Error(w, "Ошибка при преобразовании в JSON:", http.StatusInternalServerError)
		return
	}
	// fmt.Println(exp)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
func expressionsId(w http.ResponseWriter, r *http.Request) {
	parts := r.URL.Path
	parts = parts[len("/api/v1/expressions/:"):]
	id, err := strconv.Atoi(parts)
	if err != nil {
		http.Error(w, "Ошибка при конвертации id в число", http.StatusBadRequest)
		return
	}
	if m[id].Id != 0 {
		idExp := Expressions{
			Id:     id,
			Status: m[id].Status,
			Result: m[id].Result,
		}
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.MarshalIndent(idExp, "", "    ")
		if err != nil {
			http.Error(w, "Ошибка при преобразовании в JSON:", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	} else {
		http.Error(w, "По такому Id ничего не найдено", http.StatusNotFound)
		return
	}

}
