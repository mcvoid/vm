package vm

import (
	"bytes"
	"testing"
)

func TestFactorialRecursive(t *testing.T) {
	vm, _ := New(1 << 10)
	output := &bytes.Buffer{}
	vm.Stdout = output

	fact, recur := 8, 15
	n := 0

	prog := Program{
		Push, 5,
		Call, fact, 1,
		Print, 1,
		Halt,

		LoadArg, n, // fact
		JumpIfNotZero, recur,
		Push, 1,
		Return,

		LoadArg, n, // recur
		LoadArg, n,
		Push, 1,
		Subtract,
		Call, fact, 1,
		Multiply,
		Return,
	}
	vm.Run(prog, 0, false)
	if output.String() != "120\n" {
		t.Error("Factorial test error: expected fact(5) = 120, actual: ", output)
	}
}

func TestFactorialTailRecursive(t *testing.T) {
	vm, _ := New(1 << 10)
	output := &bytes.Buffer{}
	vm.Stdout = output

	fact, recur := 10, 17
	n, acc := 0, 1

	prog := Program{
		Push, 5,
		Push, 1,
		Call, fact, 2,
		Print, 1,
		Halt,

		LoadArg, n, //fact
		JumpIfNotZero, recur,
		LoadArg, acc,
		Return,

		LoadArg, n, // recur
		LoadArg, acc,
		Multiply,
		StoreArg, acc,
		LoadArg, n,
		Push, 1,
		Subtract,
		StoreArg, n,
		Push, 0,
		JumpIfZero, fact,
	}

	vm.Run(prog, 0, false)
	if output.String() != "120\n" {
		t.Error("Factorial tail resursion test error: expected fact(5) = 120, actual: ", output)
	}
}
