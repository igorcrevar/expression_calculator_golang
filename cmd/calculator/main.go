package main

import (
	pkg "expression_calculator/pkg"
	"fmt"
	"os"
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
	if len(os.Args) < 2 {
		fmt.Println("at least one expression must be passed as argument!")
	} else {
		for i := 1; i < len(os.Args); i += 1 {
			printCalculation(os.Args[i])
		}
	}
}
