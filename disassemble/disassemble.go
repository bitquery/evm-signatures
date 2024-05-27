package disassemble

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"
)

type Program struct {
	FunctionSignatures []string
	EventSignatures    []string
}

func Disassemble(code []byte) *Program {
	result := &Program{
		FunctionSignatures: make([]string, 0),
		EventSignatures:    make([]string, 0),
	}

	instructions := LoadInstructionsFromBytecode(code)

	for i := 4; i < len(instructions); {
		seq := instructions[i-4 : i+1]
		i++

		//  MAIN SOLC PATTERN
		//  DUP1 PUSHN <SELECTOR> EQ PUSH2/3 <OFFSET> JUMPI
		//  https://github.com/ethereum/solidity/blob/58811f134ac369b20c2ec1120907321edf08fff1/libsolidity/codegen/ContractCompiler.cpp#L332
		if seq[0].Eq(vm.DUP1) && (seq[1].IsPush() && seq[1].LessOrEqual(vm.PUSH4)) && seq[2].Eq(vm.EQ) && seq[3].IsPush() && seq[4].Eq(vm.JUMPI) {
			sig := hexutil.Encode(seq[1].Arg)
			result.FunctionSignatures = append(result.FunctionSignatures, sig)
		}

		continue
	}

	return result
}
