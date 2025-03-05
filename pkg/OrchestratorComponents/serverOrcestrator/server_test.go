package serverorcestrator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculate(t *testing.T) {
	expr := MathExpr{Expression: "3 + 5"}
	jsonData, _ := json.Marshal(expr)

	req, err := http.NewRequest("POST", "/api/v1/calculate/", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(calculate)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var idResp IdExpressions
	if err := json.Unmarshal(rr.Body.Bytes(), &idResp); err != nil {
		t.Fatalf("Expected valid JSON response; got: %s", rr.Body.String())
	}

	if idResp.Id <= 0 {
		t.Errorf("Expected valid ID; got: %d", idResp.Id)
	}
}

func TestExpressions(t *testing.T) {
	Ids.Id = 1
	m[1] = Expressions{Id: 1, Status: "принято", Result: ""}
	m[2] = Expressions{Id: 2, Status: "выполнено", Result: "8"}

	req, err := http.NewRequest("GET", "/api/v1/expressions/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(expressions)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var expressionsResponse []Expressions
	if err := json.Unmarshal(rr.Body.Bytes(), &expressionsResponse); err != nil {
		t.Errorf("Could not unmarshal response: %v", err)
	}
	if len(expressionsResponse) != 2 {
		t.Errorf("Expected 2 expressions; got: %d", len(expressionsResponse))
	}
}

func TestTask(t *testing.T) {
	Rtmp = make(map[int]TmpOper)
	Rtmp[1] = TmpOper{Id: 1, Num1: "3", Num2: "5", Operator: "+"}

	req, err := http.NewRequest("GET", "/internal/task/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(task)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var taskResponse Task
	if err := json.Unmarshal(rr.Body.Bytes(), &taskResponse); err != nil {
		t.Errorf("Could not unmarshal task response: %v", err)
	}
	if taskResponse.Id != 1 {
		t.Errorf("Expected task ID 1; got: %d", taskResponse.Id)
	}

	taskResp := R{Id: 1, Result: "8"}
	taskData, _ := json.Marshal(taskResp)

	reqPost, err := http.NewRequest("POST", "/internal/task/", bytes.NewBuffer(taskData))
	if err != nil {
		t.Fatal(err)
	}
	rrPost := httptest.NewRecorder()
	handlerPost := http.HandlerFunc(task)

	reqPost.Header.Set("Content-Type", "application/json")

	handlerPost.ServeHTTP(rrPost, reqPost)

	if status := rrPost.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code for POST: got %v want %v",
			status, http.StatusOK)
	}

	if _, ok := Rtmp[1]; ok {
		t.Errorf("Expected task 1 to be removed from p.Rtmp, but it's still present")
	}
}
