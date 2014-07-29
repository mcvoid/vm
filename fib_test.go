package vm

import (
    "testing"
    "bytes"
)

func TestFibonacci(t *testing.T) {
	vm, _ := New(1 << 10)
	output := &bytes.Buffer{}
	vm.Stdout = output

	fib := 8
	test_eq_1 := 15
	recur := 25
	prog := Program{
		Push, 5, //0
		Call, fib, 1, //2
		Print, 1, //5
		Halt, //7
		// start fib(n)
		// if n == 0 return 0
		LoadArg, 0, // 8, fib
		JumpIfNotZero, test_eq_1, //10
		Push, 0, //12
		Return, //14
		// if n == 1 return 1
		LoadArg, 0, // 15, test_eq_1
		Push, 1, // 17
		Subtract,             // 19
		JumpIfNotZero, recur, //20
		Push, 1, //22
		Return, //24
		// else
		// a = fib(n-1)
		// b = fib(n-2)
		// c = a + b
		// return c
		LoadArg, 0, // 25, recur
		Push, 1,
		Subtract,
		Call, fib, 1,
		LoadArg, 0,
		Push, 2,
		Subtract,
		Call, fib, 1,
		Add,
		Return,
	}
	vm.Run(prog, 0, false)
	if output.String() != "5\n" {
		t.Error("Error on fibonacci test: expected: 5\n actual: ", output.String())
	}
}

func TestFibonacciTail(t *testing.T) {
	vm, _ := New(1 << 10)
	output := &bytes.Buffer{}
	vm.Stdout = output

	fib, recur := 12, 19
	a, b, n := 0, 1, 2
	prog := Program{
		Push, 0,
		Push, 1,
		Push, 5,
		Call, fib, 3,
		Print, 1,
		Halt,
		// if n == 0 return a
		LoadArg, n, // fib
		JumpIfNotZero, recur,
		LoadArg, a,
		Return,
		// else
		LoadArg, n, //recur
		Push, 1,
		Subtract,
		LoadArg, a,
		LoadArg, b,
		Add,
		LoadArg, b,
		StoreArg, a,
		StoreArg, b,
		StoreArg, n,
		Push, 0,
		JumpIfZero, fib,
	}

	vm.Run(prog, 0, false)
	if output.String() != "5\n" {
		t.Error("Error on fibonacci test: expected: 5\n actual: ", output.String())
	}
}