package pkg

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

type funcType int64

type binaryFn func(a float64, b float64) float64

const (
	noneFunc funcType = 0
	sqrtFunc funcType = 1
)

const defaultOperator byte = '+'

type calculationStackItem struct {
	stack    Stack
	operator byte
	function funcType
}

func Calculate(expression string) (float64, error) {
	expression = strings.ToLower(expression)

	stackOfStacks := Stack{} // make(Stack, 4)
	var current calculationStackItem = calculationStackItem{stack: Stack{}, operator: defaultOperator, function: noneFunc}
	position := 0

	for position < len(expression) {
		if expression[position] >= '0' && expression[position] <= '9' {
			number, newPosition := readNumber(expression, position)
			position = newPosition
			err := executeOperation(&current, number)
			if err != nil {
				return 0, err
			}
		} else {
			ch := expression[position]
			position += 1
			switch ch {
			case '(':
				stackOfStacks.Push(current)
				current = calculationStackItem{stack: Stack{}, operator: defaultOperator, function: noneFunc}
			case ')':
				subResult := sumStack(&current.stack)
				item, isOk := stackOfStacks.Pop()
				if !isOk {
					return 0, errors.New("invalid closing parenthesis")
				}
				current = item.(calculationStackItem)
				err := executeOperation(&current, subResult)
				if err != nil {
					return 0, err
				}
			case '+', '/', '*', '^':
				current.operator = ch
			case '-':
				if current.operator == '-' {
					current.operator = '+'
				} else {
					current.operator = '-'
				}
			case ' ': // Do nothing
			default:
				// only low letters for now are
				if ch == 's' && position+4 < len(expression) && expression[position:position+4] == "qrt(" {
					current.function = sqrtFunc
					position += 3 // not 4 because in next step ( case will be executed
				} else {
					return 0, fmt.Errorf("invalid operator %c", ch)
				}
			}
		}
	}

	if !stackOfStacks.IsEmpty() {
		return 0, errors.New("invalid count of opened parenthesis")
	}

	result := sumStack(&current.stack)
	return result, nil
}

func sumStack(stack *Stack) float64 {
	var result float64 = 0
	for !stack.IsEmpty() {
		value, _ := stack.Pop()
		result += value.(float64)
	}
	return result
}

func readNumber(expression string, position int) (float64, int) {
	var number float64 = 0
	for position < len(expression) && expression[position] >= '0' && expression[position] <= '9' {
		number = number*10 + float64(expression[position]-'0')
		position += 1
	}

	if position < len(expression) && expression[position] == '.' {
		var divider float64 = 10
		position += 1
		for position < len(expression) && expression[position] >= '0' && expression[position] <= '9' {
			number += float64(expression[position]-'0') / divider
			divider *= 10
			position += 1
		}
	}

	return number, position
}

func executeOperation(current *calculationStackItem, number float64) error {
	// execute function on number if needed
	switch current.function {
	case sqrtFunc:
		number = math.Sqrt(number)
	}

	lastOperator := funcType(current.operator)
	// reset function and operator
	current.function = noneFunc
	current.operator = defaultOperator

	switch lastOperator {
	case '/':
		if number == 0 {
			return errors.New("divide by zero")
		}
		return executeBinary(current, number, func(a float64, b float64) float64 { return a / b }, "divide")
	case '*':
		return executeBinary(current, number, func(a float64, b float64) float64 { return a * b }, "multiply")
	case '^':
		return executeBinary(current, number, func(a float64, b float64) float64 { return math.Pow(a, b) }, "pow")
	case '-':
		current.stack.Push(-number)
		return nil
	default:
		current.stack.Push(number)
		return nil
	}
}

func executeBinary(current *calculationStackItem, number float64, fn binaryFn, text string) error {
	prev, isOk := current.stack.Pop()
	if !isOk {
		return fmt.Errorf("try to execute %s but no previous number for %f", text, number)
	}

	result := fn(prev.(float64), number)
	current.stack.Push(result)
	return nil
}
