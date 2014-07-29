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
	}
	for _, f := range tests {
		vm, err = New(1 << 10)
		output = &bytes.Buffer{}
		vm.Stdout = output
		f()
	}
}
