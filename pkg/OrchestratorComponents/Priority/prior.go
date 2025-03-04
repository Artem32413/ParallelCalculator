package priority

import (
	"errors"
	"fmt"
	// "fmt"
	"slices"
	"strings"
)

type R struct {
	Id     int    `json:"id"`
	Result string `json:"result"`
}
type TmpOper struct {
	Id       int    `json:"id"`
	Num1     string `json:"num1"`
	Operator string `json:"operator"`
	Num2     string `json:"num2"`
}

var RPerem = make(map[int]R)     // Тут результат от агента
var Rtmp = make(map[int]TmpOper) // Промежуточные действия

func Priority(id int, text []byte) string {
	expression := string(text)
	expression = strings.ReplaceAll(expression, " ", "")
	if err := validateExpression(expression); err != nil {
		panic(err)
	}
	var curEl string
	var l string
	var sl []string

	for _, el := range string(text) {
		curEl = string(el)
		if curEl == "*" || curEl == "/" || curEl == "-" || curEl == "+" || curEl == "(" || curEl == ")" {
			if l != "" {
				sl = append(sl, l)
			}
			sl = append(sl, curEl)
			l = ""
		} else {
			l += curEl
		}
	}

	if l != "" {
		sl = append(sl, l)
	}
	// Пробежимся по всем скобкам
	for {
		var ok bool
		sl, ok = inBracket(id, sl)
		if !ok {
			break
		}
	}
	str, _ := mainCalc(id, sl)
	// fmt.Println(str)
	return str
}
func inBracket(id int, sl []string) ([]string, bool) {
	var q []string
	var bracket1 int
	var bracket2 int
	insideBrackets := false

	for i, el := range sl {
		if el == "(" {
			bracket1 = i
			insideBrackets = true
			q = []string{}
			continue
		}
		if el == ")" {
			if !insideBrackets {
				return sl, false
			}
			str, err := mainCalc(id, q)
			if err != true {
				return sl, false
			}
			bracket2 = i
			sl[bracket1] = str

			sl = slices.Delete(sl, bracket1+1, bracket2+1)

			insideBrackets = false
			return sl, true
		}
		if insideBrackets {
			q = append(q, el)
		}
	}
	return sl, false
}
func validateExpression(expression string) error {
	lastWasDigit := false
	parensCount := 0

	for i, char := range expression {
		if (char >= '0' && char <= '9') || char == '.' {
			lastWasDigit = true
		} else {
			if lastWasDigit == false {
				if char != '(' {
					return errors.New("некорректный символ перед: " + string(char))
				}
			}
			switch char {
			case '+', '-', '*', '/':
				if i == 0 || !lastWasDigit {
					return errors.New("некорректный оператор: " + string(char))
				}
				lastWasDigit = false

			case '(':
				parensCount++

			case ')':
				parensCount--
				if parensCount < 0 {
					return errors.New("несоответствующие скобки")
				}

			default:
				return errors.New("неизвестный символ: " + string(char))
			}
		}
	}

	if parensCount != 0 {
		return errors.New("несоответствующие скобки")
	}
	// if !lastWasDigit {
	// 	return errors.New("выражение должно заканчиваться цифрой")
	// }
	return nil
}
func mainCalc(id int, k []string) (string, bool) {
	var ok bool
	for {
		k, ok = priority(id, k)
		if !ok {
			return k[0], true
		}
	}
}
func priority(id int, z []string) ([]string, bool) {
	for i, el := range z {
		if el == "*" || el == "/" {
			return Run(id, z, i), true
		}
	}
	for i, el := range z {
		if el == "+" || el == "-" {
			return Run(id, z, i), true
		}
	}
	return z, false
}
func Run(id int, z []string, i int) []string {
	setTmpOper(id, z[i-1], z[i], z[i+1])
	// Ждать решения

	for {
		if _, ok := RPerem[id]; ok {
			z[i-1] = RPerem[id].Result
			delete(RPerem, id)
			fmt.Println("result", z[i-1])
			break
		}
	}
	// fmt.Println(RPerem)
	d := slices.Delete(z, i, i+2)
	return d
}
func setTmpOper(id int, num1, i, num2 string) { // Добавляем действие в очередь для вычисления
	num1 = strings.TrimSpace(num1)
	i = strings.TrimSpace(i)
	num2 = strings.TrimSpace(num2)
	fmt.Println(num1, i, num2)
	Rtmp[id] = TmpOper{
		Id:       id,
		Num1:     num1,
		Operator: i,
		Num2:     num2,
	}
}
