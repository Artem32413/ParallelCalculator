package calculation

import (
	"fmt"
	"log"
	"strconv"
)

func Calcularion(n1, sign, n2 string) string {
	
	in1, err := strconv.ParseFloat(n1, 64)
	fmt.Println(in1)
	if err != nil {
		fmt.Println(err)
		log.Printf("Ошибка конвертации 1 числа\n")
		return ""
	}
	in2, err := strconv.ParseFloat(n2, 64)
	if err != nil {
		log.Printf("Ошибка конвертации 2 числа\n")
		return ""
	}
	switch sign {
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
