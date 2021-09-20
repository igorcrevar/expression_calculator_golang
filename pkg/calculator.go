package pkg

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

type funcType int64

const (
	noneFunc funcType = 0
	sqrtFunc funcType = 1
)

const (
	defaultOperator byte = ' '
	skipOperator    byte = 'S'
)

type binaryFn func(a float64, b float64) float64

type calculationItem struct {
	operator byte
	number   float64
}

type subProcessItem struct {
	operations  []calculationItem
	operator    byte
	modificator byte
	function    funcType
}

func Calculate(expression string) (float64, error) {
	expression = strings.ToLower(expression)

	stackOfStacks := Stack{} // make(Stack, 4)
	var current subProcessItem = subProcessItem{
		operations: []calculationItem{}, function: noneFunc, operator: defaultOperator, modificator: 0,
	}

	var position int = 0
	var number float64 = 0.0

	for position < len(expression) {
		if expression[position] >= '0' && expression[position] <= '9' {
			number, position = readNumber(expression, position)
			err := addOperation(&current, number)
			if err != nil {
				return 0, err
			}
		} else {
			ch := expression[position]
			position += 1
			switch ch {
			case '(':
				stackOfStacks.Push(current)
				current = subProcessItem{
					operations: []calculationItem{}, function: noneFunc, operator: defaultOperator, modificator: 0,
				}
			case ')':
				item, isOk := stackOfStacks.Pop()
				if !isOk {
					return 0, errors.New("invalid closing parenthesis")
				}
				number, err := calculateOperations(&current.operations)
				if err != nil {
					return 0, err
				}
				current = item.(subProcessItem)
				err = addOperation(&current, number)
				if err != nil {
					return 0, err
				}
			case '/', '*', '^':
				if current.operator != defaultOperator {
					return 0, fmt.Errorf("operator %c already specified before %c", current.operator, ch)
				}
				current.operator = ch
			case '+':
				if current.operator == defaultOperator {
					current.operator = '+'
				}
			case '-':
				if current.operator == defaultOperator {
					current.operator = '-'
				} else if current.modificator != '-' {
					current.modificator = '-'
				} else {
					current.modificator = 0
				}
			case ' ': // Do nothing
			default:
				// only low letters for now are
				if isFunc, newPosition := isFunction(expression, position-1, &current); isFunc {
					position = newPosition
				} else {
					return 0, fmt.Errorf("invalid operator %c", ch)
				}
			}
		}
	}

	if !stackOfStacks.IsEmpty() {
		return 0, errors.New("invalid count of opened parenthesis")
	}

	return calculateOperations(&current.operations)
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

func addOperation(processItem *subProcessItem, number float64) error {
	var err error = nil
	// number must be preceded with some operator if its not first number
	if len(processItem.operations) > 0 && processItem.operator == defaultOperator {
		err = fmt.Errorf("operator not specified for %f", number)
	} else {
		number = executeFunction(processItem, number)
		if processItem.modificator == '-' {
			number = -number
		}

		processItem.operations = append(processItem.operations, calculationItem{
			number: number, operator: processItem.operator,
		})
		processItem.modificator = defaultOperator
		processItem.operator = defaultOperator
	}
	return err
}

func calculateOperations(operations *[]calculationItem) (float64, error) {
	var err error = nil

	// first priority ^ from right to left
	for i := len(*operations) - 1; i >= 0; {
		switch (*operations)[i].operator {
		case '^':
			i, err = executeBinary(operations, i,
				func(a float64, b float64) float64 { return math.Pow(a, b) }, "pow")
		default:
			i -= 1
		}
		if err != nil {
			return 0, err
		}
	}

	// second priority from left to right
	for i := 0; i < len(*operations); i += 1 {
		switch (*operations)[i].operator {
		case '*':
			_, err = executeBinary(operations, i,
				func(a float64, b float64) float64 { return a * b }, "multiply")
		case '/':
			if (*operations)[i].number != 0 {
				_, err = executeBinary(operations, i,
					func(a float64, b float64) float64 { return a / b }, "divide")
			} else {
				err = errors.New("divide by zero")
			}
		}
		if err != nil {
			return 0, err
		}
	}

	// last priority +, -(a-b is just a+(-b))
	result := 0.0
	for _, item := range *operations {
		if item.operator == '-' {
			result -= item.number
		} else {
			result += item.number
		}
	}

	return result, nil
}

func executeFunction(current *subProcessItem, number float64) float64 {
	// execute function on number if needed
	switch current.function {
	case sqrtFunc:
		current.function = noneFunc
		return math.Sqrt(number)
	default:
		return number
	}
}

func isFunction(expression string, position int, processItem *subProcessItem) (bool, int) {
	if position+5 < len(expression) && expression[position:position+5] == "sqrt(" {
		processItem.function = sqrtFunc
		position += 4 // not 5 because in next step ( case will be executed
	}
	return processItem.function != noneFunc, position
}

func executeBinary(operations *[]calculationItem, index int, fn binaryFn, text string) (int, error) {
	second := (*operations)[index]
	firstIndex := index - 1
	for firstIndex >= 0 && (*operations)[firstIndex].operator == skipOperator {
		firstIndex -= 1
	}
	if firstIndex < 0 {
		return firstIndex, fmt.Errorf("try to execute %s but no previous number for %f", text, second.number)
	}

	first := (*operations)[firstIndex]
	result := fn(first.number, second.number)
	(*operations)[firstIndex] = calculationItem{operator: first.operator, number: result}
	(*operations)[index] = calculationItem{operator: skipOperator, number: 0}
	return firstIndex, nil
}
