package pkg

import (
	"fmt"
	"testing"
)

func TestSimple(t *testing.T) {
	result, err := Calculate("1 + - 2.1")
	if err != nil {
		t.Errorf("got error %q", err)
	} else if !almostEqual(result, -1.1) {
		t.Errorf("got result %f instead of %f", result, -1.1)
	}
}

func TestModificators(t *testing.T) {
	result, err := Calculate("-5+-3*-3*-3/-3*-3++3-log(3*-3*1^0*-1)")
	if err != nil {
		t.Errorf("got error %q", err)
	} else if !almostEqual(result, -29.95424250944) {
		t.Errorf("got result %f instead of %f", result, -29.95424250944)
	}
}

func TestParenthesisComplex(t *testing.T) {
	result, err := Calculate("1.89 + 2 * (4+(3-1)*(3*(-3+9))) -8.49")
	if err != nil {
		t.Errorf("got error %q", err)
	} else if !almostEqual(result, 73.4) {
		t.Errorf("got result %f instead of %f", result, 73.4)
	}
}

func TestMinusModificator(t *testing.T) {
	result, err := Calculate("-3^-3*---3")
	if err != nil {
		t.Errorf("got error %q", err)
	} else if !almostEqual(result, 0.11111111111) {
		t.Errorf("got result %f instead of %f", result, 0.11111111111)
	}
}

func TestSqrtAndPow(t *testing.T) {
	result, err := Calculate("-sqrt(3+3*(2+2)) +2^(3+(1+1)+1)")
	if err != nil {
		t.Errorf("got error %q", err)
	} else if !almostEqual(result, 60.1270166538) {
		t.Errorf("got result %f instead of %f", result, 60.1270166538)
	}
}

func TestParenthesisSimple(t *testing.T) {
	result, err := Calculate("((-3)*2)")
	if err != nil {
		t.Errorf("got error %q", err)
	} else if !almostEqual(result, -6.0) {
		t.Errorf("got result %f instead of %f", result, -6.0)
	}
}

func TestDivideByZero(t *testing.T) {
	_, err := Calculate("4/(1-1)")
	if err == nil || err.Error() != "divide by zero" {
		t.Errorf("got error %q instead of %q", err, "divide by zero")
	}
}

func TestInvalidParenthesisOpen(t *testing.T) {
	_, err := Calculate("(4/(1+1)")
	if err == nil || err.Error() != "invalid count of opened parenthesis" {
		t.Errorf("got error %q instead of %q", err, "invalid count of opened parenthesis")
	}
}

func TestInvalidParenthesisClose(t *testing.T) {
	_, err := Calculate("(4/(4-1--1))))")
	if err == nil || err.Error() != "invalid closing parenthesis" {
		t.Errorf("got error %q instead of %q", err, "invalid closing parenthesis")
	}
}

func TestInvalidDivide(t *testing.T) {
	should := "try to execute divide but no previous number for 3.000000"
	_, err := Calculate("3+(/3)")
	if err == nil || err.Error() != should {
		t.Errorf("got error %q instead of %q", err, should)
	}
}

func TestInvalidMultiply(t *testing.T) {
	should := "try to execute multiply but no previous number for 3.000000"
	_, err := Calculate("*3")
	if err == nil || err.Error() != should {
		t.Errorf("got error %q instead of %q", err, should)
	}
}

func TestInvalidPow(t *testing.T) {
	should := "try to execute pow but no previous number for 3.000000"
	_, err := Calculate("^3")
	if err == nil || err.Error() != should {
		t.Errorf("got error %q instead of %q", err, should)
	}
}

func TestInvalidOperator(t *testing.T) {
	should := "invalid operator @"
	_, err := Calculate("3@-6")
	if err == nil || err.Error() != should {
		t.Errorf("got error %q instead of %q", err, should)
	}
}

func TestOperatorNotSpecified(t *testing.T) {
	should := "operator not specified for 5.000000"
	_, err := Calculate("4 - 4 5")
	if err == nil || err.Error() != should {
		t.Errorf("got error %q instead of %q", err, should)
	}
}

func TestOperatorNotSpecifiedParenthesis(t *testing.T) {
	should := "operator not specified for -1024.000000"
	_, err := Calculate("4 (-4^5)")
	if err == nil || err.Error() != should {
		t.Errorf("got error %q instead of %q", err, should)
	}
}

func TestOperatorAlreadySpecified(t *testing.T) {
	for _, ch := range [3]byte{'/', '^', '*'} {
		should := fmt.Sprintf("operator * already specified before %c", ch)
		_, err := Calculate(fmt.Sprintf("4*%c3", ch))
		if err == nil || err.Error() != should {
			t.Errorf("got error %q instead of %q", err, should)
		}
	}
}

func almostEqual(a, b float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}

	return diff <= 0.000001
}
