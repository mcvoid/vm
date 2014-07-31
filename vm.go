/*
Package vm provides a program-embeddable virtual machine which can
execute scripts which act upon io.Reader and io.Writer.

Features
- A hot-swappable instruction set.
- STDIN, STDOUT, STDERR readable and writable to any io.Reader or io.Writer
- Adjustable memory size
- Enable scripting into your Go app.

How to use:
- import "github.com/mcvoid/vm"
- vm, err := vm.New(1024)
- prog := vm.Program{vm.Push, 3, vm.Push, 5, vm.Add, vm.Print, 1, vm.Halt}
- vm.Run(prog, 0, false)
- Watch it print "8" to vm.Stdout (os.Stdout by default)

What it can't do:
- compile a higher-level language to machine code (maybe in a separate package)
- Stack bounds checking (Go has bounds checking, though, so it will just panic on overflow)
- Printing ASCII/UTF-8 (yet)
- Reading from Stdin (yet)
- Self-modifying code (code and stack reside in different areas)
- Multi-word return values (maybe in the future)
- Not even remotely thread-safe. Use different VM's to execute code concurrently.

What it can do:
- Function calls.
- Fibonacci numbers! Ackermann functions! Factorials!
- Tail recursion.
- Add your own instructions to make it do more.
- Hook vm.Stdout to an http.ResponseWriter to script your web apps!

*/
package vm

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// VM represents the state of a virtual machine.
// Mem and the memory in which code resides are separate spaces.
type VM struct {
	Mem    []int          // Mem represents the  internal Memory of the VM. It does not include code.
	SP     int            // SP, or the stack pointer, is the address of the top of the stack. It points to the next available unused space.
	FP     int            // FP, or the frame pointer, is the address of the start of the current function's stack frame.
	IP     int            // IP, or the instruction pointer, is the address of the next instruction to execute.
	IS     InstructionSet // The set of instructions which can be executed on this VM.
	Stdin  io.Reader      // Stdin is the VM's standard input reader. This is the reader which Read reads from.
	Stdout io.Writer      // Stdout is the VM's standard output writer. This is what Print writes to.
	Stderr io.Writer      // Stderr is the VM's standard error writer. This is what stack traces write to.
}

// A program is a sequence of integers.
// Instructions, offsets, and literal values are all integers.
type Program []int

// A fake io.Writer which discards everything written to it, like /dev/null.
// Useful for suppressing a VM's output on Stdout or Stderr.
type Bitbucket struct{}
func (bit Bitbucket)Write(b []byte) (n int, err error) { return len(b), nil }

// New creates a new virtual machine.
//
// size is the size of the VM's stack memory in words. This does not include the memory used by code, so available stack size is unaffected
// by program size. The size of a word is the same size as an int.
func New(size int) (vm *VM, err error) {
	if size > (1<<32-1) || size < 0 {
		return vm, errors.New("vm: size not within range")
	}
	vm = &VM{
		Mem:   make([]int, size),
		IS:    Default,
		Stdin: os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	return vm, err
}

// Runs a program on the VM
// src is the source of the bytecode
// start is the index of the code starting point
func (vm *VM) Run(src Program, start int) {
	vm.IP = start
	// run the program until it halts
	for {
		if vm.IP < 0 || vm.IP >= len(src) {
			fmt.Fprintln(vm.Stderr, "IP ", vm.IP, " out of bounds.")
			break
		}
		instr := src[vm.IP]
		if vm.IP+1+vm.IS[instr].Args > len(src) {
			fmt.Fprintf(vm.Stderr, "IP out of bounds on arguments of ", vm.IS[instr].Text)
			break
		}
		args := src[vm.IP+1 : vm.IP+1+vm.IS[instr].Args]
		fmt.Fprintln(vm.Stderr, vm.IP, ":\t", vm.IS[instr].Text, args, "\tstack:\t", vm.Mem[:vm.SP])
		vm.IP += 1 + len(args)
		vm.IS[instr].Action(vm, args)
		if vm.IS[instr].Halts {
			break
		}
	}
}

// Push pushes a value to the stack, updating SP accordingly and in the correct sequence.
func (vm *VM) Push(val int) {
	vm.Mem[vm.SP] = val
	vm.SP++
}

// Pop pops a value from the stack, updating SP accordingly and in the correct sequence.
func (vm *VM) Pop() int {
	vm.SP--
	val := vm.Mem[vm.SP]
	return val
}
