package pkg

import (
	"errors"
	"fmt"
)

type calculationStackItem struct {
	stack    Stack
	operator byte
}

func Calculate(expression string) (float64, error) {
	stackOfStacks := Stack{} // make(Stack, 4)
	var current calculationStackItem = calculationStackItem{stack: Stack{}, operator: '+'}
	position := 0

	for position < len(expression) {
		if expression[position] >= '0' && expression[position] <= '9' {
			number, newPosition := readNumber(expression, position)
			position = newPosition
			err := executeOperation(&current, number)
			if err != nil {
				return 0, err
			}
			current.operator = '+'
		} else {
			ch := expression[position]
			position += 1
			switch ch {
			case '(':
				stackOfStacks.Push(current)
				current = calculationStackItem{stack: Stack{}, operator: '+'}
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
			case '+', '/', '*':
				current.operator = ch
			case '-':
				if current.operator == '-' {
					current.operator = '+'
				} else {
					current.operator = '-'
				}
			case ' ': // Do nothing
			default:
				return 0, fmt.Errorf("invalid operator %c", ch)
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
	switch current.operator {
	case '/':
		prev, isOk := current.stack.Pop()
		if !isOk {
			return fmt.Errorf("try to divide but no previous number for %f", number)
		} else if number == 0 {
			return errors.New("divide by zero")
		}

		current.stack.Push(prev.(float64) / number)
		return nil
	case '*':
		prev, isOk := current.stack.Pop()
		if !isOk {
			return fmt.Errorf("try to multiply but no previous number for %f", number)
		}

		current.stack.Push(prev.(float64) * number)
		return nil
	case '-':
		current.stack.Push(-number)
		return nil
	default:
		current.stack.Push(number)
		return nil
	}
}
