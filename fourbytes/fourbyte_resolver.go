package fourbytes

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"

	signatures "github.com/bitquery/evm-signatures"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

const (
	MethodIdLength = 4
	EventIdLength  = 32
)

//go:embed events.json
var events []byte

// go:embed methods.json
var methods []byte

type FourByteResolver struct {
	methodsDB map[string]string
	eventsDB  map[string]string
}

func NewFourByteResolver(_, eventsPath string) (*FourByteResolver, error) {
	resolver := &FourByteResolver{
		methodsDB: make(map[string]string),
		eventsDB:  make(map[string]string),
	}

	if err := json.Unmarshal(methods, &resolver.methodsDB); err != nil {
		return nil, fmt.Errorf("failed to read methods db: %v", err)
	}

	if err := json.Unmarshal(events, &resolver.eventsDB); err != nil {
		return nil, fmt.Errorf("failed to read events db: %v", err)
	}

	return resolver, nil
}

func (f *FourByteResolver) ResolveContractABI(code []byte) (*abi.ABI, error) {
	sigs := signatures.FindContractSignatures(code)
	result := abi.ABI{
		Methods: make(map[string]abi.Method),
		Events:  make(map[string]abi.Event),
	}

	for _, s := range sigs.FunctionSignatures {
		sel, err := f.MethodSelector(s)
		if err != nil {
			return nil, err
		}

		selector, err := abi.ParseSelector(sel)
		if err != nil {
			return nil, err
		}

		args, err := convertArgs(selector.Inputs)
		if err != nil {
			return nil, err
		}

		method := abi.NewMethod(selector.Name, sel, abi.Function, "", false, false, args, nil)
		result.Methods[method.Name] = method
	}

	for _, s := range sigs.EventSignatures {
		sel, err := f.EventSelector(s)
		if err != nil {
			return nil, err
		}

		selector, err := abi.ParseSelector(sel)
		if err != nil {
			return nil, err
		}

		args, err := convertArgs(selector.Inputs)
		if err != nil {
			return nil, err
		}

		event := abi.NewEvent(selector.Name, sel, false, args)
		result.Events[event.Name] = event
	}

	return &result, nil
}

// Selector checks the given 4byte ID against the known ABI methods.
//
// This method does not validate the match, it's assumed the caller will do.
func (f *FourByteResolver) MethodSelector(id []byte) (string, error) {
	if len(id) < MethodIdLength {
		return "", fmt.Errorf("expected 4-byte id")
	}

	sig := hex.EncodeToString(id[:MethodIdLength])
	selector, exists := f.methodsDB[sig]
	if !exists {
		return selector, fmt.Errorf("method signature not found")
	}

	return selector, nil
}

func (f *FourByteResolver) EventSelector(id []byte) (string, error) {
	if len(id) < EventIdLength {
		return "", fmt.Errorf("expected 32-byte id")
	}

	sig := hex.EncodeToString(id[:EventIdLength])
	selector, exists := f.eventsDB[sig]
	if !exists {
		return selector, fmt.Errorf("event signature not found")
	}

	return selector, nil
}

func convertArgs(args []abi.ArgumentMarshaling) (abi.Arguments, error) {
	r := make(abi.Arguments, 0, len(args))
	for _, arg := range args {
		t, err := abi.NewType(arg.Type, arg.InternalType, arg.Components)
		if err != nil {
			return nil, fmt.Errorf("failed to convert args: %v", err)
		}

		r = append(r, abi.Argument{
			Name:    arg.Name,
			Type:    t,
			Indexed: arg.Indexed,
		})
	}

	return r, nil
}
