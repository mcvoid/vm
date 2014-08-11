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
	   Load reads a value from Memory address addr and Pushes it on the stack.
	   syntax: load <addr>
	*/
	Load
	/*
		Store pops a value from the stack and writes it to Memory at address addr.
		syntax: store <addr>
	*/
	Store
	/*
	   LoadArg Pushes the nth argument to the stack.
	   syntax: loadarg <n>
	*/
	LoadArg
	/*
	   StorArg pops a value from the stack and stores it in the nth argument.
	   syntax: storarg <n>
	*/
	StoreArg
	/*
	   Add pops two values from the stack and Pushes the sum of those values.
	   syntax: add
	*/
	Add
	/*
	   Substract pops two values from the stack and Pushes the difference of the second value popped minus the first.
	   syntax: sub
	*/
	Subtract
	/*
	   Multiply pops two values from the stack and Pushes the product of the two..
	   syntax: mul
	*/
	Multiply
	/*
	   Divide pops two values from the stack and returns the quotient of the two values minus the remainder.
	   syntax: div
	*/
	Divide
	/*
	   Modulo pops two values from the stack, divides them, and returns the quotient of the two values.
	   syntax: mod
	*/
	Modulo
	/*
	   And pops two values from the stack and returns the bitwise-and of them.
	   syntax: and
	*/
	And
	/*
	   Or pops two values from the stack and returns the bitwise-or of them.
	   syntax: or
	*/
	Or
	/*
	   Xor pops two values from the stack and returns the bitwise-xor of them.
	   syntax: xor
	*/
	Xor
	/*
	   Not pops two values from the stack and returns the bitwise-not of them.
	   syntax: not
	*/
	Not
	/*
		ShiftLeft pops two values from the stack and returns the first value, shifted
		to the left a number of bits as the second value.
		syntax: shl
	*/
	ShiftLeft
	/*
		ShiftRight pops two values from the stack and returns the first value,
		logically shifted to the left a number of bits as the second value.
		syntax: shr
	*/
	ShiftRight
	/*
	   Pushes a literal value to the stack.
	   syntax: Push <val>
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
	   Pops a value from the current stack frame, returns from the function, and Pushes the value back.
	   syntax: ret
	*/
	Return
	/*
	   Prints n values to the console.
	   syntax: print <n>
	*/
	Print
	//Scan
)

// Instruction contains critical information about a single instruction for the VM.
type Instruction struct {
	// Action is the code which is called that is executed upon the vm to change its state.
	// args is any arguments to the instruction which follow in the bytecode.
	// len(args) = Args
	Action func(vm *VM, args []int)
	Text   string // Text is the textual representaion of the instruction as seen in stack traces.
	Args   int    // Args tells the VM how many arguments to pull from the code.
	Halts  bool   // Halts determines whether the instruction halts execution after running this instruction.
}

// Represents an instruction in string form.
func (i Instruction) String() string {
	return i.Text
}

// InstructionSet is a collection of all of the instructions the VM can execute in a given program.
type InstructionSet []Instruction

// The behavior of each bytecode
var Default InstructionSet

func init() {
	Default = InstructionSet{
		Instruction{
			func(vm *VM, args []int) {},
			"halt",
			0,
			true,
		},
		Instruction{
			func(vm *VM, args []int) {
				addr := args[0]
				val := vm.Mem[addr]
				vm.Push(val)
			},
			"load",
			1,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				addr := args[0]
				val := vm.Pop()
				vm.Mem[addr] = val
			},
			"store",
			1,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				numargs := vm.Mem[vm.FP-3]
				addr := vm.FP - 3 - numargs + args[0]
				val := vm.Mem[addr]
				vm.Push(val)
			},
			"loadarg",
			1,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				numargs := vm.Mem[vm.FP-3]
				addr := vm.FP - 3 - numargs + args[0]
				val := vm.Pop()
				vm.Mem[addr] = val
			},
			"storarg",
			1,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				val2, val1 := vm.Pop(), vm.Pop()
				vm.Push(val1 + val2)
			},
			"add",
			0,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				val2, val1 := vm.Pop(), vm.Pop()
				vm.Push(val1 - val2)
			},
			"sub",
			0,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				val2, val1 := vm.Pop(), vm.Pop()
				vm.Push(val1 * val2)
			},
			"mul",
			0,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				val2, val1 := vm.Pop(), vm.Pop()
				vm.Push(val1 / val2)
			},
			"div",
			0,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				val2, val1 := vm.Pop(), vm.Pop()
				vm.Push(val1 % val2)
			},
			"mod",
			0,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				val2, val1 := vm.Pop(), vm.Pop()
				vm.Push(val1 & val2)
			},
			"and",
			0,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				val2, val1 := vm.Pop(), vm.Pop()
				vm.Push(val1 | val2)
			},
			"or",
			0,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				val2, val1 := vm.Pop(), vm.Pop()
				vm.Push(val1 ^ val2)
			},
			"xor",
			0,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				val1 := vm.Pop()
				vm.Push(^val1)
			},
			"not",
			0,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				val2, val1 := vm.Pop(), vm.Pop()
				vm.Push(val1 << uint(val2))
			},
			"shl",
			0,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				val2, val1 := vm.Pop(), vm.Pop()
				vm.Push(val1 >> uint(val2))
			},
			"shr",
			0,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				val := args[0]
				vm.Push(val)
			},
			"Push",
			1,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				vm.Pop()
			},
			"pop",
			0,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				addr := args[0]
				if vm.Pop() == 0 {
					vm.IP = addr
				}
			},
			"jz",
			1,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				addr := args[0]
				if vm.Pop() != 0 {
					vm.IP = addr
				}
			},
			"jnz",
			1,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				addr, numargs := args[0], args[1]
				vm.Push(numargs) // save number of arguments (for popping later)
				vm.Push(vm.FP)   // save old frame pointer
				vm.Push(vm.IP)   // save return address
				vm.FP = vm.SP    // make new stack frame at return address
				vm.IP = addr     // jump!
			},
			"call",
			2,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				val := vm.Pop()
				vm.SP = vm.FP    // restore the stack
				addr := vm.Pop() // get return address
				vm.FP = vm.Pop() // get previous frame pointer
				numargs := vm.Pop()
				for i := 0; i < numargs; i++ {
					vm.Pop()
				}
				vm.Push(val)
				vm.IP = addr // jump!
			},
			"ret",
			0,
			false,
		},
		Instruction{
			func(vm *VM, args []int) {
				n := args[0]
				for i := 0; i < n; i++ {
					val := vm.Pop()
					fmt.Fprint(vm.Stdout, val)
				}
				fmt.Fprintln(vm.Stdout)
                vm.Stdout.Flush()
			},
			"print",
			1,
			false,
		},
	}
}
