package vm

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// a virtual CPU and its memory
type VM struct {
	mem    []int
	code   Program
	sp     int
	fp     int
	ip     int
	IS     InstructionSet
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type Program []int

// Creates a new virtual machine
// size is the size of the VM's memory in words
func NewVM(size int) (vm *VM, err error) {
	if size > (1<<32-1) || size < 0 {
		return vm, errors.New("vm: size not within range")
	}
	mem := make([]int, size)
	vm = &VM{sp: 0, ip: 0, fp: 0}
	vm.mem = mem
	vm.IS = Default
	vm.Stdin = os.Stdin
	vm.Stdout = os.Stdout
	vm.Stderr = os.Stderr
	return vm, err
}

// Runs a program on the VM
// src is the source of the bytecode
// start is the index of the code starting point
// main is the word index of the program entry point
func (vm *VM) Run(src Program, start int, trace bool) {
	vm.code = src
	vm.ip = start
	// run the program until it halts
	for {
		if vm.ip < 0 || vm.ip >= len(vm.code) {
			fmt.Fprintln(vm.Stderr, "IP ", vm.ip, " out of bounds.")
			break
		}
		instr := vm.code[vm.ip]
		if vm.ip+1+vm.IS[instr].Args > len(vm.code) {
			fmt.Fprintf(vm.Stderr, "IP out of bounds on arguments of ", vm.IS[instr].Text)
			break
		}
		args := vm.code[vm.ip+1 : vm.ip+1+vm.IS[instr].Args]
		if trace {
			fmt.Fprintln(vm.Stderr, vm.ip, ":\t", vm.IS[instr].Text, args, "\tstack:\t", vm.mem[:vm.sp])
		}
		vm.ip += 1 + len(args)
		vm.IS[instr].Action(vm, args)
		if vm.IS[instr].Halts {
			break
		}
	}
}

// pushes a value on a stack
func (vm *VM) push(val int) {
	vm.mem[vm.sp] = val
	vm.sp++
}

// pops a value from the stack
func (vm *VM) pop() int {
	vm.sp--
	val := vm.mem[vm.sp]
	return val
}
