package vm

import (
	"fmt"
)

// Instructions built-in to the Default instruction set.
const (
	/*
	   Halt tells the VM to stop executing the current program.
	   syntax: halt
	*/
	Halt int = iota
	/*
	   Load reads a value from memory address addr and pushes it on the stack.
	   syntax: load <addr>
	*/
	Load
	/*
		Store pops a value from the stack and writes it to memory at address addr.
		syntax: store <addr>
	*/
	Store
	/*
	   LoadArg pushes the nth argument to the stack.
	   syntax: loadarg <n>
	*/
	LoadArg
	/*
	   StorArg pops a value from the stack and stores it in the nth argument.
	   syntax: storarg <n>
	*/
	StoreArg
	/*
	   Add pops two values from the stack and pushes the sum of those values.
	   syntax: add
	*/
	Add
	/*
	   Substract pops two values from the stack and pushes the difference of the seccond value popped minus the first.
	   syntax: sub
	*/
	Subtract
	/*
	   Pushes a literal value to the stack.
	   syntax: push <val>
	*/
	Push
	/*
	   Pops a value from the stack, discarding it.
	   syntax: pop
	*/
	Pop
	/*
	   JumpIfZero Pops a value from the stack. If the value is zero, it branches to addr.
	   syntax: jz <addr>
	*/
	JumpIfZero
	/*
	   JumpIfNotZero pops a value from the stack. If the value is not zero, it branches to addr.
	   syntax: jnz <addr>
	*/
	JumpIfNotZero
	/*
	   Calls a function with n arguments at address addr.
	   syntax: call <addr> <n>
	*/
	Call
	/*
	   Pops a value from the current stack frame, returns from the function, and pushes the value back.
	   syntax: ret
	*/
	Return
	/*
	   Prints n values to the console.
	   syntax: print <n>
	*/
	Print
)

// Critical information about an instruction
type Instruction struct {
	Action func(*VM, []int) // code that is executed by the instruction.
	Text   string           // the textual representaion of the instruction.
	Args   int              // the number of arguments it has.
	Halts  bool             // Whether the instruction halts execution.
}

func (i Instruction) String() string {
	return i.Text
}

// All of the instructions the VM can
type InstructionSet []Instruction

// The behavior of each bytecode
var Default InstructionSet = InstructionSet{
	Instruction{
		func(vm *VM, args []int) {},
		"halt",
		0,
		true,
	},
	Instruction{
		func(vm *VM, args []int) {
			addr := args[0]
			val := vm.mem[addr]
			vm.push(val)
		},
		"load",
		1,
		false,
	},
	Instruction{
		func(vm *VM, args []int) {
			addr := args[0]
			val := vm.pop()
			vm.mem[addr] = val
		},
		"store",
		1,
		false,
	},
	Instruction{
		func(vm *VM, args []int) {
			numargs := vm.mem[vm.fp-3]
			addr := vm.fp - 3 - numargs + args[0]
			val := vm.mem[addr]
			vm.push(val)
		},
		"loadarg",
		1,
		false,
	},
	Instruction{
		func(vm *VM, args []int) {
			numargs := vm.mem[vm.fp-3]
			addr := vm.fp - 3 - numargs + args[0]
			val := vm.pop()
			vm.mem[addr] = val
		},
		"storarg",
		1,
		false,
	},
	Instruction{
		func(vm *VM, args []int) {
			val2, val1 := vm.pop(), vm.pop()
			vm.push(val1 + val2)
		},
		"add",
		0,
		false,
	},
	Instruction{
		func(vm *VM, args []int) {
			val2, val1 := vm.pop(), vm.pop()
			vm.push(val1 - val2)
		},
		"sub",
		0,
		false,
	},
	Instruction{
		func(vm *VM, args []int) {
			val := args[0]
			vm.push(val)
		},
		"push",
		1,
		false,
	},
	Instruction{
		func(vm *VM, args []int) {
			vm.pop()
		},
		"pop",
		0,
		false,
	},
	Instruction{
		func(vm *VM, args []int) {
			addr := args[0]
			if vm.pop() == 0 {
				vm.ip = addr
			}
		},
		"jz",
		1,
		false,
	},
	Instruction{
		func(vm *VM, args []int) {
			addr := args[0]
			if vm.pop() != 0 {
				vm.ip = addr
			}
		},
		"jnz",
		1,
		false,
	},
	Instruction{
		func(vm *VM, args []int) {
			addr, numargs := args[0], args[1]
			vm.push(numargs) // save number of arguments (for popping later)
			vm.push(vm.fp)   // save old frame pointer
			vm.push(vm.ip)   // save return address
			vm.fp = vm.sp    // make new stack frame at return address
			vm.ip = addr     // jump!
		},
		"call",
		2,
		false,
	},
	Instruction{
		func(vm *VM, args []int) {
			val := vm.pop()
			vm.sp = vm.fp    // restore the stack
			addr := vm.pop() // get return address
			vm.fp = vm.pop() // get previous frame pointer
			numargs := vm.pop()
			for i := 0; i < numargs; i++ {
				vm.pop()
			}
			vm.push(val)
			vm.ip = addr // jump!
		},
		"ret",
		0,
		false,
	},
	Instruction{
		func(vm *VM, args []int) {
			n := args[0]
			for i := 0; i < n; i++ {
				val := vm.pop()
				fmt.Fprint(vm.Stdout, val)
			}
			fmt.Fprintln(vm.Stdout)
		},
		"print",
		1,
		false,
	},
}
