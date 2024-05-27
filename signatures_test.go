package signatures

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed testdata/USDT.bin
var USDTBytecode []byte

//go:embed testdata/BAYC.bin
var BAYCBytecode []byte

func Test_Disassemble(t *testing.T) {
	t.Parallel()

	ERC20DefaultFunctionSignatures := []string{
		"0xdd62ed3e", // allowance(address,address)
		"0x095ea7b3", // approve(address,uint256)
		"0x70a08231", // balanceOf(address)
		"0x18160ddd", // totalSupply()
		"0xa9059cbb", // transfer(address,uint256)
		"0x23b872dd", // transferFrom(address,address,uint256)
	}
	ERC721DefaultFunctionSignatures := []string{
		"0x70a08231", // balanceOf(address)
		"0x6352211e", // ownerOf(uint256)
		"0xb88d4fde", // safeTransferFrom(address,address,uint256,bytes)
		"0x42842e0e", // safeTransferFrom(address,address,uint256)
		"0x23b872dd", // transferFrom(address,address,uint256)
		"0x095ea7b3", // approve(address,uint256)
		"0xa22cb465", // setApprovalForAll(address,bool)
		"0x081812fc", // getApproved(uint256),
		"0xe985e9c5", // isApprovedForAll(address,address)
	}
	tests := []struct {
		name                string
		code                []byte
		wantFuncSignatures  []string
		wantEventSignatures []string
	}{
		{
			name:               "erc20/USDT",
			code:               USDTBytecode,
			wantFuncSignatures: ERC20DefaultFunctionSignatures,
			wantEventSignatures: []string{
				"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", // Transfer(address,address,uint256)
			},
		},
		{
			name:               "erc721/BAYC",
			code:               BAYCBytecode,
			wantFuncSignatures: ERC721DefaultFunctionSignatures,
			wantEventSignatures: []string{
				"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", // Transfer(address,address,uint256)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := FindContractSignatures(tt.code)

			for _, w := range tt.wantFuncSignatures {
				assert.Contains(t, got.FunctionSignatures, w)
			}

			for _, w := range tt.wantEventSignatures {
				assert.Contains(t, got.EventSignatures, w)
			}
		})
	}
}
