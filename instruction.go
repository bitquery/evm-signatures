package evm_signatures

import (
	"github.com/ethereum/go-ethereum/core/asm"
	"github.com/ethereum/go-ethereum/core/vm"
)

type Instruction struct {
	OpCode vm.OpCode
	Arg    []byte
	PC     uint64
}

func (i *Instruction) IsLog() bool {
	return i.OpCode >= vm.LOG0 && i.OpCode <= vm.LOG4
}

func (i *Instruction) IsHalt() bool {
	return i.OpCode == vm.STOP || i.OpCode == vm.RETURN || i.OpCode >= vm.REVERT
}

func (i *Instruction) IsPush() bool {
	return i.OpCode >= vm.PUSH1 && i.OpCode <= vm.PUSH32
}

func (i *Instruction) Eq(op vm.OpCode) bool {
	return i.OpCode == op
}

func (i *Instruction) LessThan(op vm.OpCode) bool {
	return i.OpCode < op
}

func (i *Instruction) LessOrEqual(op vm.OpCode) bool {
	return i.OpCode <= op
}

func LoadInstructionsFromBytecode(code []byte) []*Instruction {
	it := asm.NewInstructionIterator(code)
	var instructions []*Instruction
	for it.Next() {
		instructions = append(instructions, &Instruction{
			OpCode: it.Op(),
			Arg:    it.Arg(),
			PC:     it.PC(),
		})
	}
	return instructions
}
