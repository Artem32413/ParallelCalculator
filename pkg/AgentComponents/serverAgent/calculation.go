package serverAgent

import (
	"fmt"
	"log"
	"strconv"
)

func Calcularion(miniTask Task) string {

	in1, err := strconv.ParseFloat(miniTask.Arg1, 64)
	fmt.Println(in1)
	if err != nil {
		log.Printf("Ошибка конвертации 1 числа\n")
		return ""
	}
	in2, err := strconv.ParseFloat(miniTask.Arg2, 64)
	if err != nil {
		log.Printf("Ошибка конвертации 2 числа\n")
		return ""
	}
	switch miniTask.Operation {
	case "+":
		res := in1 + in2
		result := strconv.FormatFloat(res, 'f', 2, 64)
		return result
	case "-":
		res := in1 - in2
		result := strconv.FormatFloat(res, 'f', 2, 64)
		return result
	case "*":
		res := in1 * in2
		result := strconv.FormatFloat(res, 'f', 2, 64)
		return result
	case "/":
		if in2 == 0 {
			log.Println("Деление на ноль")
			return ""
		}
		res := in1 / in2
		result := strconv.FormatFloat(res, 'f', 2, 64)
		return result
	default:
		log.Println("Неизвестная операция")
		return ""
	}
}
func ParalelCalc(miniTask Task) {
	eq := Calcularion(miniTask)
	var res = Result{
		Id:     miniTask.Id,
		Result: eq,
	}
	PostResult(res)
}
