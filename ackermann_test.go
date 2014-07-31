package vm

import (
	"bytes"
	"testing"
)

// Implements the Ackermann Phi function,
// Which is a generalization of addition, multiplication,
// exponentiation, and all hyperoperations
func TestAckermannPhi(t *testing.T) {
	vm, _ := New(1 << 10)
	output := &bytes.Buffer{}
	vm.Stdout = output

	phi, case2, case3, case4, recur := 12, 22, 36, 46, 49
	a, b, n := 0, 1, 2

	prog := Program{
		Push, 3,
		Push, 3,
		Push, 2,
		Call, phi, 3,
		Print, 1,
		Halt,

		// base case 1: a+b if n=0
		LoadArg, n, // line 12
		JumpIfNotZero, case2,
		LoadArg, a,
		LoadArg, b,
		Add,
		Return,

		// base case 2, 0 if b=0 && n=1
		LoadArg, b, //line 22
		JumpIfNotZero, recur,
		LoadArg, n,
		Push, 1,
		Subtract,
		JumpIfNotZero, case3,
		Push, 0,
		Return,

		// base case 3, 1 if b=0 && n=2
		LoadArg, n, // line 36
		Push, 2,
		Subtract,
		JumpIfNotZero, case4,
		Push, 1,
		Return,

		// base case 4, a if b=0 && n>2
		LoadArg, a, // line 46
		Return,

		// recursive case: Phi(a, Phi(a, b-1, n), n-1)
		LoadArg, a, // line 49
		LoadArg, a,
		LoadArg, b,
		Push, 1,
		Subtract,
		LoadArg, n,
		Call, phi, 3,
		LoadArg, n,
		Push, 1,
		Subtract,
		Call, phi, 3,
		Return,
	}

	vm.Run(prog, 0, false)
	if output.String() != "27\n" {
		t.Error("Error on Ackermann Phi test: Phi(3, 3, 2) = 27, actual: ", output.String())
	}
}
