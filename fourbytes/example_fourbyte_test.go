package fourbytes_test

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/bitquery/evm-signatures/fourbytes"
)

//go:embed SEAPORT.bin
var seaport []byte

func Example() {
	r, err := fourbytes.NewFourByteResolver("methods.json", "events.json")
	if err != nil {
		panic(err)
	}

	res, err := r.ResolveContractABI(seaport)
	if err != nil {
		panic(err)
	}
	m, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(m))
}
