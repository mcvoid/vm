VM - An simple Virtual Machine for Go.

PACKAGE DOCUMENTATION

package vm
    import "github.com/mcvoid/vm"

    Package vm provides a program-embeddable virtual machine which can
    execute scripts which act upon io.Reader and io.Writer.

    Features - A hot-swappable instruction set. - STDIN, STDOUT, STDERR
    readable and writable to any io.Reader or io.Writer - Adjustable memory
    size - Enable scripting into your Go app.

    How to use: - import "github.com/mcvoid/vm" - vm, err := vm.New(1024) -
    prog := vm.Program{vm.Push, 3, vm.Push, 5, vm.Add, vm.Print, 1, vm.Halt}
    - vm.Run(prog, 0) - Watch it print "8" to vm.Stdout (os.Stdout by
    default)

    What it can't do: - compile a higher-level language to machine code
    (maybe in a separate package) - Stack bounds checking (Go has bounds
    checking, though, so it will just panic on overflow) - Printing
    ASCII/UTF-8 (yet) - Reading from Stdin (yet) - Self-modifying code (code
    and stack reside in different areas) - Multi-word return values (maybe
    in the future) - Not even remotely thread-safe. Use different VM's to
    execute code concurrently.

    What it can do: - Function calls. - Fibonacci numbers! Ackermann
    functions! Factorials! - Tail recursion. - Add your own instructions to
    make it do more. - Hook vm.Stdout to an http.ResponseWriter to script
    your web apps!

CONSTANTS

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
       Substract pops two values from the stack and pushes the difference of the second value popped minus the first.
       syntax: sub
    */
    Subtract
    /*
       Multiply pops two values from the stack and pushes the product of the two..
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
    Instructions built-in to the Default instruction set.

TYPES

type Instruction struct {
    // Action is the code which is called that is executed upon the vm to change its state.
    // args is any arguments to the instruction which follow in the bytecode.
    // len(args) = Args
    Action func(vm *VM, args []int)
    Text   string // Text is the textual representaion of the instruction as seen in stack traces.
    Args   int    // Args tells the VM how many arguments to pull from the code.
    Halts  bool   // Halts determines whether the instruction halts execution after running this instruction.
}
    Instruction contains critical information about a single instruction for
    the VM.

func (i Instruction) String() string
    Represents an instruction in string form.

type InstructionSet []Instruction
    InstructionSet is a collection of all of the instructions the VM can
    execute in a given program.

var Default InstructionSet
    The behavior of each bytecode

type Program []int
    A program is a sequence of integers. Instructions, offsets, and literal
    values are all integers.

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
    VM represents the state of a virtual machine. Mem and the memory in
    which code resides are separate spaces.

func New(size int) (vm *VM, err error)
    New creates a new virtual machine.

    size is the size of the VM's stack memory in words. This does not
    include the memory used by code, so available stack size is unaffected
    by program size. The size of a word is the same size as an int.

func (vm *VM) Pop() int
    Pop pops a value from the stack, updating SP accordingly and in the
    correct sequence.

func (vm *VM) Push(val int)
    Push pushes a value to the stack, updating SP accordingly and in the
    correct sequence.

func (vm *VM) Run(src Program, start int)
    Runs a program on the VM src is the source of the bytecode start is the
    index of the code starting point


