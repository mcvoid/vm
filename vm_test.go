package vm

import (
	"bytes"
	"testing"
	_ "time"
)

func testSmoke(vm *VM, t *testing.T) {
	prog := Program{Halt}
	vm.Run(prog, 0, false)
}

func testPush(vm *VM, t *testing.T) {
	prog := Program{
		Push, 10,
		Halt,
	}
	vm.Run(prog, 0, false)
	if vm.mem[0] != 10 {
		t.Error("Error on push test: Expected ", []int{10}, " Actual ", vm.mem[:1])
	}
}

func testPrint(vm *VM, t *testing.T) {
	prog := Program{
		Push, 10,
		Print, 1,
		Halt,
	}
	vm.Run(prog, 0, false)
	if vm.Stdout.(*bytes.Buffer).String() != "10\n" {
		t.Error("Error on Print test: expected: 10\n actual: ", vm.Stdout.(*bytes.Buffer).String())
	}
}

func testPop(vm *VM, t *testing.T) {
	prog := Program{
		Push, 10,
		Push, 20,
		Pop,
		Print, 1,
		Halt,
	}
	vm.Run(prog, 0, false)
	if vm.Stdout.(*bytes.Buffer).String() != "10\n" {
		t.Error("Error on pop test: expected: 10\n actual: ", vm.Stdout.(*bytes.Buffer).String())
	}
}

func testStore(vm *VM, t *testing.T) {
	prog := Program{
		Push, 10,
		Store, 100,
		Halt,
	}
	vm.Run(prog, 0, false)
	if vm.mem[100] != 10 {
		t.Error("Error on stor test: vm.mem[100] expected 10 actual: ", vm.mem[100])
	}
}

func testLoad(vm *VM, t *testing.T) {
	prog := Program{
		Load, 100,
		Print, 1,
		Halt,
	}
	vm.mem[100] = 10
	vm.Run(prog, 0, false)
	if vm.Stdout.(*bytes.Buffer).String() != "10\n" {
		t.Error("Error on load test: expected: 10\n actual: ", vm.Stdout.(*bytes.Buffer).String())
	}
}

func testAdd(vm *VM, t *testing.T) {
	prog := Program{
		Push, 10,
		Push, 20,
		Add,
		Print, 1,
		Halt,
	}
	vm.Run(prog, 0, false)
	if vm.Stdout.(*bytes.Buffer).String() != "30\n" {
		t.Error("Error on add test: expected: 30\n actual: ", vm.Stdout.(*bytes.Buffer).String())
	}
}

func testSubtract(vm *VM, t *testing.T) {
	prog := Program{
		Push, 10,
		Push, 20,
		Subtract,
		Print, 1,
		Halt,
	}
	vm.Run(prog, 0, false)
	if vm.Stdout.(*bytes.Buffer).String() != "-10\n" {
		t.Error("Error on sub test: expected: -10\n actual: ", vm.Stdout.(*bytes.Buffer).String())
	}
}

func testJumpIfNotZero(vm *VM, t *testing.T) {
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
	if vm.Stdout.(*bytes.Buffer).String() != "20\n10\n" {
		t.Error("Error on jnz test: expected: 20\n10\n actual: ", vm.Stdout.(*bytes.Buffer).String())
	}
}

func testJumpIfZero(vm *VM, t *testing.T) {
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
	if vm.Stdout.(*bytes.Buffer).String() != "10\n20\n" {
		t.Error("Error on jz test: expected: 10\n20\n actual: ", vm.Stdout.(*bytes.Buffer).String())
	}
}

func testCallReturn(vm *VM, t *testing.T) {
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
	if vm.Stdout.(*bytes.Buffer).String() != "4\n" {
		t.Error("Error on call/return test: expected: 4\n actual: ", vm.Stdout.(*bytes.Buffer).String())
	}
}

func testLoadArg(vm *VM, t *testing.T) {
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
	if vm.Stdout.(*bytes.Buffer).String() != "5\n" {
		t.Error("Error on loadarg 1-parameter test: expected: 5\n actual: ", vm.Stdout.(*bytes.Buffer).String())
	}
	vm.Stdout.(*bytes.Buffer).Truncate(0)
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
	if vm.Stdout.(*bytes.Buffer).String() != "5\n" {
		t.Error("Error on loadarg 2-parameter test: expected: 5\n actual: ", vm.Stdout.(*bytes.Buffer).String())
	}
}

func testStoreArg(vm *VM, t *testing.T) {
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
	if vm.Stdout.(*bytes.Buffer).String() != "3\n" {
		t.Error("Error on storarg test: expected: 3\n actual: ", vm.Stdout.(*bytes.Buffer).String())
	}
}

func TestRun(t *testing.T) {
	tests := []func(*VM, *testing.T){testSmoke, testPush, testPrint, testPop, testStore,
		testLoad, testAdd, testSubtract, testJumpIfNotZero, testJumpIfZero,
		testCallReturn, testLoadArg, testStoreArg,
	}

	vm, err := New(1 << 10)
	if err != nil {
		t.Error("Error initializing VM")
	}
	for _, f := range tests {
		vm, err = New(1 << 10)
		output := &bytes.Buffer{}
		vm.Stdout = output
		f(vm, t)
	}
}
