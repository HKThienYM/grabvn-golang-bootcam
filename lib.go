package main

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

func isDigit(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func isOpt(char rune) bool {
	if char == '+' || char == '-' || char == '*' || char == '/' {
		return true
	}
	return false
}

func deleteBlank(s []string) []string {
	var r []string
	for _, str := range s {
		if str != " " && str != "" {
			r = append(r, str)
		}
	}
	return r
}

func validateExpresion(exprsn string) error {
	for i, char := range exprsn {
		if !((char >= '0' && char <= '9') || isOpt(char) || char == ' ' || char == '(' || char == ')') {
			return errors.New("Error: Wrong format")
		}
		if isOpt(char) && isOpt(rune(exprsn[i+1])) {
			return errors.New("Error: Wrong format")
		}
	}
	input := strings.Split(exprsn, " ")
	input = deleteBlank(input)
	if len(input) == 0 {
		return errors.New("Error: No input value")
	}
	if input[0][0] == '*' || input[0][0] == '/' {
		return errors.New("Error: Wrong format: " + input[0])
	}
	if isOpt(rune(input[len(input)-1][len(input[len(input)-1])-1])) {
		return errors.New("Error: Wrong format: " + input[len(input)-1])
	}
	for i := 0; i < len(input); i++ {
		if isDigit(string(input[i][len(input[i])-1])) && i < len(input)-1 {
			if isDigit(string(input[i+1][0])) {
				return errors.New("Error: Wrong format: " + input[i] + " " + input[i+1])
			}
		}
		if isOpt(rune(input[i][len(input[i])-1])) && i < len(input)-1 {
			if isOpt(rune(input[i+1][0])) {
				return errors.New("Error: Wrong format: " + input[i] + " " + input[i+1])
			}
		}
	}
	return nil
}

func calculateString(exprsn string) (result float64, err error) {
	err = validateExpresion(exprsn)
	if err != nil {
		return
	}
	prefix := infixToPrefix(exprsn)
	result, err = evaluatePrefix(prefix)
	return
}

func getPriority(opt string) int {
	if opt == "*" || opt == "/" {
		return 2
	} else if opt == "+" || opt == "-" {
		return 1
	} else {
		return 0
	}
}

func infixToPrefix(infix string) string {
	var operators, operands []string

	for i := 0; i < len(infix); i++ {
		if infix[i] == ' ' {
			continue
		} else if infix[i] == '(' {
			operators = append(operators, "(")
		} else if infix[i] == ')' {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				op1 := operands[len(operands)-1]
				op2 := operands[len(operands)-2]
				operands = operands[:len(operands)-2]

				op := operators[len(operators)-1]
				operators = operators[:len(operators)-1]

				operands = append(operands, op+" "+op2+" "+op1)
			}
			operators = operators[:len(operators)-1]
		} else if isDigit(string(infix[i])) {
			var tmp string
			for ; i < len(infix) && infix[i] >= '0' && infix[i] < '9'; i++ {
				tmp = tmp + string(infix[i])
			}
			operands = append(operands, tmp)
			i--
			continue

		} else {
			if len(operators) == 0 && len(operands) == 0 && (infix[i] == '+' || infix[i] == '-') {
				operands = append(operands, "0")
			}

			for len(operators) > 0 && getPriority(string(infix[i])) <= getPriority(operators[len(operators)-1]) {

				op1 := operands[len(operands)-1]
				op2 := operands[len(operands)-2]
				operands = operands[:len(operands)-2]

				op := operators[len(operators)-1]
				operators = operators[:len(operators)-1]

				operands = append(operands, op+" "+op2+" "+op1)
			}

			operators = append(operators, string(infix[i]))
		}
	}
	for len(operators) > 0 {

		op1 := operands[len(operands)-1]
		op2 := operands[len(operands)-2]
		operands = operands[:len(operands)-2]

		op := operators[len(operators)-1]
		operators = operators[:len(operators)-1]

		operands = append(operands, op+" "+op2+" "+op1)
	}
	return operands[0]
}

func evaluatePrefix(exprsn string) (float64, error) {
	input := strings.Split(exprsn, " ")
	var stack []float64

	for j := len(input) - 1; j >= 0; j-- {
		if isDigit(input[j]) {
			//controlled
			f, _ := strconv.ParseFloat(input[j], 64)
			stack = append(stack, f)
		} else {
			o1 := stack[len(stack)-1]
			o2 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			switch input[j] {
			case "+":
				stack = append(stack, o1+o2)
			case "-":
				stack = append(stack, o1-o2)
			case "*":
				stack = append(stack, o1*o2)
			case "/":
				if o2 == 0 {
					return 0, errors.New("Error: Can't divide by zero")
				}
				stack = append(stack, o1/o2)

			}
		}
	}

	return stack[0], nil
}
