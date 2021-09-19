package pkg

import "testing"

func TestSimple(t *testing.T) {
	result, err := Calculate("1 + - 2.1")
	if err != nil {
		t.Errorf("got error %q", err)
	}
	if !almostEqual(result, -1.1) {
		t.Errorf("got result %f instead of %f", result, -1.1)
	}
}

func TestParenthesisComplex(t *testing.T) {
	result, err := Calculate("1.89 + 2 * (4+(3-1)*(3*(-3+9))) -8.49")
	if err != nil {
		t.Errorf("got error %q", err)
	}
	if !almostEqual(result, 73.4) {
		t.Errorf("got result %f instead of %f", result, 73.4)
	}
}

func TestParenthesisSimple(t *testing.T) {
	result, err := Calculate("((-3)*2)")
	if err != nil {
		t.Errorf("got error %q", err)
	}
	if !almostEqual(result, -6.0) {
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
	should := "try to divide but no previous number for 6.000000"
	_, err := Calculate("/(3+3)")
	if err == nil || err.Error() != should {
		t.Errorf("got error %q instead of %q", err, should)
	}
}

func TestInvalidMultiply(t *testing.T) {
	should := "try to multiply but no previous number for 3.000000"
	_, err := Calculate("*3")
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

func almostEqual(a, b float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}

	return diff <= 0.000001
}
