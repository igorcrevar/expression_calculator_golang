package main

import (
	pkg "expression_calculator/pkg"
	"fmt"
)

func printCalculation(expression string) {
	result, err := pkg.Calculate(expression)
	if err == nil {
		fmt.Printf("%s = %.4f\n", expression, result)
	} else {
		fmt.Printf("%s\n", err)
	}
}

func main() {
	printCalculation("1 + 2 * (4+(3-1)*(3*(-3+9))) -8")
	printCalculation("1 + 7 * 4")
	printCalculation("4/(1-1)")
}
