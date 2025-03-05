package serverAgent

import (
	"testing"
)

func TestCalcularion(t *testing.T) {
	testCases := []struct {
		task     Task
		expected string
	}{
		{task: Task{Arg1: "2", Operation: "+", Arg2: "2"}, expected: "4.00"},
		{task: Task{Arg1: "5", Operation: "-", Arg2: "3"}, expected: "2.00"},
		{task: Task{Arg1: "4", Operation: "*", Arg2: "6"}, expected: "24.00"},
		{task: Task{Arg1: "10", Operation: "/", Arg2: "2"}, expected: "5.00"},
		{task: Task{Arg1: "7.5", Operation: "+", Arg2: "2.5"}, expected: "10.00"},
		{task: Task{Arg1: "10", Operation: "-", Arg2: "2.5"}, expected: "7.50"},
		{task: Task{Arg1: "3.2", Operation: "*", Arg2: "2"}, expected: "6.40"},
		{task: Task{Arg1: "5", Operation: "/", Arg2: "2"}, expected: "2.50"},
		{task: Task{Arg1: "10", Operation: "/", Arg2: "0"}, expected: ""},
		{task: Task{Arg1: "5", Operation: "%", Arg2: "2"}, expected: ""},
		{task: Task{Arg1: "abc", Operation: "+", Arg2: "2"}, expected: ""},
		{task: Task{Arg1: "5", Operation: "+", Arg2: "xyz"}, expected: ""},
		{task: Task{Arg1: "-5", Operation: "+", Arg2: "2"}, expected: "-3.00"},
		{task: Task{Arg1: "5", Operation: "+", Arg2: "-2"}, expected: "3.00"},
		{task: Task{Arg1: "0", Operation: "*", Arg2: "100"}, expected: "0.00"},
	}

	for _, tc := range testCases {
		t.Run(tc.task.Arg1+tc.task.Operation+tc.task.Arg2, func(t *testing.T) {
			actual := Calcularion(tc.task)
			if actual != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, actual)
			}
		})
	}
}
