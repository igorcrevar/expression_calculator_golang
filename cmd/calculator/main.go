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
		fmt.Printf("%s = %s\n", expression, err)
	}
}

func main() {
	printCalculation("1+log(3*-3*1^0*-1)")
	printCalculation("-5^2 + -5^2")
	printCalculation("1+2-3^2")
	printCalculation("4^(1-1)+sqrt(4^2)")
	printCalculation("4/(1-1)")
	printCalculation("1 + 2 * (4+(3-1)*(3*(-3+9))) -8")
	printCalculation("1 + 7 * 4")
	printCalculation("4 - 4 5")
}
