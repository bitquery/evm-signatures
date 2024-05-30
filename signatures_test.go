package signatures_test

import (
	_ "embed"
	"testing"

	signatures "github.com/bitquery/evm-signatures"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/USDT.bin
var USDTBytecode []byte

//go:embed testdata/BAYC.bin
var BAYCBytecode []byte

func Test_FindContractSignatures(t *testing.T) {
	t.Parallel()

	ERC20DefaultFunctionSignatures := [][]byte{
		hexutil.MustDecode("0xdd62ed3e"), // allowance(address,address)
		hexutil.MustDecode("0x095ea7b3"), // approve(address,uint256)
		hexutil.MustDecode("0x70a08231"), // balanceOf(address)
		hexutil.MustDecode("0x18160ddd"), // totalSupply()
		hexutil.MustDecode("0xa9059cbb"), // transfer(address,uint256)
		hexutil.MustDecode("0x23b872dd"), // transferFrom(address,address,uint256)
	}
	ERC721DefaultFunctionSignatures := [][]byte{
		hexutil.MustDecode("0x70a08231"), // balanceOf(address)
		hexutil.MustDecode("0x6352211e"), // ownerOf(uint256)
		hexutil.MustDecode("0xb88d4fde"), // safeTransferFrom(address,address,uint256,bytes)
		hexutil.MustDecode("0x42842e0e"), // safeTransferFrom(address,address,uint256)
		hexutil.MustDecode("0x23b872dd"), // transferFrom(address,address,uint256)
		hexutil.MustDecode("0x095ea7b3"), // approve(address,uint256)
		hexutil.MustDecode("0xa22cb465"), // setApprovalForAll(address,bool)
		hexutil.MustDecode("0x081812fc"), // getApproved(uint256),
		hexutil.MustDecode("0xe985e9c5"), // isApprovedForAll(address,address)
	}
	tests := []struct {
		name                string
		code                []byte
		wantFuncSignatures  [][]byte
		wantEventSignatures [][]byte
	}{
		{
			name:               "erc20/USDT",
			code:               USDTBytecode,
			wantFuncSignatures: ERC20DefaultFunctionSignatures,
			wantEventSignatures: [][]byte{
				hexutil.MustDecode("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"), // Transfer(address,address,uint256)
			},
		},
		{
			name:               "erc721/BAYC",
			code:               BAYCBytecode,
			wantFuncSignatures: ERC721DefaultFunctionSignatures,
			wantEventSignatures: [][]byte{
				hexutil.MustDecode("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"), // Transfer(address,address,uint256)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			got := signatures.FindContractSignatures(tt.code)

			for _, w := range tt.wantFuncSignatures {
				assert.Contains(t, got.FunctionSignatures, w)
			}

			for _, w := range tt.wantEventSignatures {
				assert.Contains(t, got.EventSignatures, w)
			}
		})
	}
}
