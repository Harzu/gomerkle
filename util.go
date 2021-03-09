package gomerkle

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

func sha3Encode(data ...interface{}) string {
	return crypto.Keccak256Hash([]byte(fmt.Sprintf("%v", data))).String()
}
