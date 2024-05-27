package evm_signatures

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"
)

type Signatures struct {
	FunctionSignatures []string
	EventSignatures    []string
}

func FindContractSignatures(code []byte) *Signatures {
	result := &Signatures{
		FunctionSignatures: make([]string, 0),
		EventSignatures:    make([]string, 0),
	}

	lastPush32Value := make([]byte, 0, 32)
	instructions := LoadInstructionsFromBytecode(code)
	seq := make([]*Instruction, 0, 5)
	for i := 0; i < len(instructions); i++ {
		inst := instructions[i]
		if inst.Eq(vm.PUSH32) {
			lastPush32Value = inst.Arg
		}

		if inst.IsLog() && len(lastPush32Value) > 0 {
			eventSig := hexutil.Encode(lastPush32Value)
			result.EventSignatures = append(result.EventSignatures, eventSig)
		}

		// take the last 5 instructions and check for function signature
		j := min(len(instructions), i+5)
		seq = instructions[i:j]
		if sig := findFunctionSignature(seq); sig != "" {
			result.FunctionSignatures = append(result.FunctionSignatures, sig)
		}
	}

	return result
}

func findFunctionSignature(seq []*Instruction) string {
	if len(seq) < 5 {
		return ""
	}

	//  MAIN SOLC PATTERN
	//  DUP1 PUSHN <SELECTOR> EQ PUSH2/3 <OFFSET> JUMPI
	//  https://github.com/ethereum/solidity/blob/58811f134ac369b20c2ec1120907321edf08fff1/libsolidity/codegen/ContractCompiler.cpp#L332
	if seq[0].Eq(vm.DUP1) && (seq[1].IsPush() && seq[1].LessOrEqual(vm.PUSH4)) && seq[2].Eq(vm.EQ) && seq[3].IsPush() && seq[4].Eq(vm.JUMPI) {
		return hexutil.Encode(seq[1].Arg)
	}

	// VYPER
	// Vyper compiler has more patterns, so we're trying to cover the most popular, but this is not a 100%

	// VYPER with XOR
	// PUSHN <SELECTOR> DUP2 XOR PUSH2 <OFFSET> JUMPI
	if (seq[0].IsPush() && seq[0].LessOrEqual(vm.PUSH4)) && seq[1].Eq(vm.DUP2) && seq[2].Eq(vm.XOR) && seq[3].IsPush() && seq[4].Eq(vm.JUMPI) {
		return hexutil.Encode(seq[0].Arg)
	}

	// VYPER with MLOAD [old versions]
	// PUSHN <SELECTOR> PUSH1 00 MLOAD EQ ISZERO JUMPI
	if (seq[0].IsPush() && seq[0].LessOrEqual(vm.PUSH4)) && seq[1].Eq(vm.PUSH1) && seq[2].Eq(vm.MLOAD) && seq[3].Eq(vm.EQ) && seq[4].Eq(vm.ISZERO) {
		return hexutil.Encode(seq[0].Arg)
	}

	return ""
}
