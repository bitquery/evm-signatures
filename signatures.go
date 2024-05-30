package signatures

import (
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	PUSH32DataLength          = 32
	InstructionSequenceLength = 5
)

type Signatures struct {
	FunctionSignatures [][]byte
	EventSignatures    [][]byte
}

func FindContractSignatures(code []byte) *Signatures {
	result := &Signatures{
		FunctionSignatures: make([][]byte, 0),
		EventSignatures:    make([][]byte, 0),
	}

	lastPush32Value := make([]byte, 0, PUSH32DataLength)
	instructions := LoadInstructionsFromBytecode(code)

	for i := 0; i < len(instructions); i++ {
		inst := instructions[i]
		if inst.Eq(vm.PUSH32) {
			lastPush32Value = inst.Arg
		}

		if inst.IsLog() && len(lastPush32Value) > 0 {
			result.EventSignatures = append(result.EventSignatures, lastPush32Value)
		}

		// take the last 5 instructions and check for function signature
		j := min(len(instructions), i+InstructionSequenceLength)
		if sig := findFunctionSignature(instructions[i:j]); sig != nil {
			result.FunctionSignatures = append(result.FunctionSignatures, sig)
		}
	}

	return result
}

func findFunctionSignature(seq []Instruction) []byte {
	if len(seq) < InstructionSequenceLength {
		return nil
	}

	//  MAIN SOLC PATTERN
	//  DUP1 PUSHN <SELECTOR> EQ PUSH2/3 <OFFSET> JUMPI
	//  https://github.com/ethereum/solidity/blob/58811f134ac369b20c2ec1120907321edf08fff1/libsolidity/codegen/ContractCompiler.cpp#L332
	if seq[0].Eq(vm.DUP1) && (seq[1].IsPush() && seq[1].LessOrEqual(vm.PUSH4)) && seq[2].Eq(vm.EQ) && seq[3].IsPush() && seq[4].Eq(vm.JUMPI) {
		return seq[1].Arg
	}

	// VYPER
	// Vyper compiler has more patterns, so we're trying to cover the most popular, but this is not a 100%

	// VYPER with XOR
	// PUSHN <SELECTOR> DUP2 XOR PUSH2 <OFFSET> JUMPI
	if (seq[0].IsPush() && seq[0].LessOrEqual(vm.PUSH4)) && seq[1].Eq(vm.DUP2) && seq[2].Eq(vm.XOR) && seq[3].IsPush() && seq[4].Eq(vm.JUMPI) {
		return seq[0].Arg
	}

	// VYPER with MLOAD [old versions]
	// PUSHN <SELECTOR> PUSH1 00 MLOAD EQ ISZERO JUMPI
	if (seq[0].IsPush() && seq[0].LessOrEqual(vm.PUSH4)) && seq[1].Eq(vm.PUSH1) && seq[2].Eq(vm.MLOAD) && seq[3].Eq(vm.EQ) && seq[4].Eq(vm.ISZERO) {
		return seq[0].Arg
	}

	return nil
}
