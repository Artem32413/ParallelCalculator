package serverorcestrator

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func ValidateExpression(expression string) error {

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
	return nil
}
func TimeSleep(i string) int {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	addition := os.Getenv("TIME_ADDITION_MS")
	subtraction := os.Getenv("TIME_SUBTRACTION_MS")
	multiplication := os.Getenv("TIME_MULTIPLICATIONS_MS")
	division := os.Getenv("TIME_DIVISIONS_MS")
	switch i {
	case "+":
		plus, err := strconv.Atoi(addition)
		if err != nil {
			panic(err)
		}
		return plus
	case "-":
		minus, err := strconv.Atoi(subtraction)
		if err != nil {
			panic(err)
		}
		return minus
	case "*":
		multiply, err := strconv.Atoi(multiplication)
		if err != nil {
			panic(err)
		}
		return multiply
	case "/":
		divide, err := strconv.Atoi(division)
		if err != nil {
			panic(err)
		}
		return divide
	default:
		return 0
	}
}
