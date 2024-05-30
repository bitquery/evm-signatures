package signatures

import "github.com/ethereum/go-ethereum/accounts/abi"

type ABIResolver interface {
	ResolveContractABI(code []byte) (*abi.ABI, error)
}
