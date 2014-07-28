package vm

import (
	"bytes"
	"testing"
	_ "time"
)

func TestRun(t *testing.T) {
	vm, err := New(1 << 10)
	output := &bytes.Buffer{}
	vm.Stdout = output
	tests := []func(){
		// making sure it's initialized
		func() {
			if err != nil {
				t.Error("Error initializing VM")
			}
		},
		// smoke test
		func() {
			prog := Program{Halt}
			vm.Run(prog, 0, false)
		},
		// test that push works
		func() {
			prog := Program{
				Push, 10,
				Halt,
			}
			vm.Run(prog, 0, false)
			if vm.mem[0] != 10 {
				t.Error("Error on push test: Expected ", []int{10}, " Actual ", vm.mem[:1])
			}
		},
		// Test that print works
		func() {
			prog := Program{
				Push, 10,
				Print, 1,
				Halt,
			}
			vm.Run(prog, 0, false)
			if output.String() != "10\n" {
				t.Error("Error on Print test: expected: 10\n actual: ", output.String())
			}
		},
		// test that pop works
		func() {
			prog := Program{
				Push, 10,
				Push, 20,
				Pop,
				Print, 1,
				Halt,
			}
			vm.Run(prog, 0, false)
			if output.String() != "10\n" {
				t.Error("Error on pop test: expected: 10\n actual: ", output.String())
			}
		},
		// test that store works
		func() {
			prog := Program{
				Push, 10,
				Store, 100,
				Halt,
			}
			vm.Run(prog, 0, false)
			if vm.mem[100] != 10 {
				t.Error("Error on stor test: vm.mem[100] expected 10 actual: ", vm.mem[100])
			}
		},
		// test that load works
		func() {
			prog := Program{
				Load, 100,
				Print, 1,
				Halt,
			}
			vm.mem[100] = 10
			vm.Run(prog, 0, false)
			if output.String() != "10\n" {
				t.Error("Error on load test: expected: 10\n actual: ", output.String())
			}
		},
		// test that add works
		func() {
			prog := Program{
				Push, 10,
				Push, 20,
				Add,
				Print, 1,
				Halt,
			}
			vm.Run(prog, 0, false)
			if output.String() != "30\n" {
				t.Error("Error on add test: expected: 30\n actual: ", output.String())
			}
		},
		// test that subtract works
		func() {
			prog := Program{
				Push, 10,
				Push, 20,
				Subtract,
				Print, 1,
				Halt,
			}
			vm.Run(prog, 0, false)
			if output.String() != "-10\n" {
				t.Error("Error on sub test: expected: -10\n actual: ", output.String())
			}
		},
		// test that jnz works
		func() {
			prog := Program{
				Push, 10,
				JumpIfNotZero, 9, // should jump
				Push, 10,
				Print, 1,
				Halt,
				Push, 20,
				Print, 1,
				Halt,
			}
			vm.Run(prog, 0, false)
			prog = Program{
				Push, 0,
				JumpIfNotZero, 9, // should not jump
				Push, 10,
				Print, 1,
				Halt,
				Push, 20,
				Print, 1,
				Halt,
			}
			vm.Run(prog, 0, false)
			if output.String() != "20\n10\n" {
				t.Error("Error on jnz test: expected: 20\n10\n actual: ", output.String())
			}
		},
		// test that jz works
		func() {
			prog := Program{
				Push, 10,
				Print, 1,
				Halt,
				Push, 0,
				JumpIfZero, 0, // should jump
				Push, 20,
				Print, 1,
				Halt,
			}
			vm.Run(prog, 5, false)
			prog = Program{
				Push, 10,
				Print, 1,
				Halt,
				Push, 10,
				JumpIfZero, 0, // should not jump
				Push, 20,
				Print, 1,
				Halt,
			}
            vm.Run(prog, 5, false)
			if output.String() != "10\n20\n" {
				t.Error("Error on jz test: expected: 10\n20\n actual: ", output.String())
			}
		},
		// test that call/return work
		func() {
			f := 6
			prog := Program{
				Call, f, 0,
				Print, 1,
				Halt,
				Push, 2,
				Push, 2,
				Add,
				Return,
			}
			vm.Run(prog, 0, false)
			if output.String() != "4\n" {
				t.Error("Error on call/return test: expected: 4\n actual: ", output.String())
			}
		},
		// test that loadarg works
		func() {
			f := 8
			prog := Program{ // 1 argument
				Push, 2,
				Call, f, 1,
				Print, 1,
				Halt,
				LoadArg, 0,
				Push, 3,
				Add,
				Return,
			}
			vm.Run(prog, 0, false)
			if output.String() != "5\n" {
				t.Error("Error on loadarg 1-parameter test: expected: 5\n actual: ", output.String())
			}
            output.Truncate(0)
			f = 10
			prog = Program{
				Push, 2,
				Push, 3,
				Call, f, 2,
				Print, 1,
				Halt,
				LoadArg, 0,
				LoadArg, 1,
				Add,
				Return,
			}
			vm.Run(prog, 0, false)
			if output.String() != "5\n" {
				t.Error("Error on loadarg 2-parameter test: expected: 5\n actual: ", output.String())
			}
		},
		// test that storarg works
		func() {
			f := 8
			prog := Program{
				Push, 2,
				Call, f, 1,
				Print, 1,
				Halt,
				Push, 3,
				StoreArg, 0,
				LoadArg, 0,
				Return,
			}
			vm.Run(prog, 0, false)
			if output.String() != "3\n" {
				t.Error("Error on storarg test: expected: 3\n actual: ", output.String())
			}
		},
        // Calculate the fibonacci sequence
        func() {
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
                Push, 1,  // 17
                Subtract, // 19
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
            vm.Run(prog, 0, true)
            if output.String() != "5\n" {
                t.Error("Error on fibonacci test: expected: 5\n actual: ", output.String())
            }
        },
        // iterative fibonacci method
        func() {
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
        },
	}
	for _, f := range tests {
        vm, err = New(1 << 10)
        output = &bytes.Buffer{}
        vm.Stdout = output
		f()
	}
}
