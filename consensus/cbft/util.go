package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
	"math/rand"
)

// Returns a random offset between 0 and n
func randomOffset(n int) int {
	if n == 0 {
		return 0
	}
	return int(rand.Uint32() % uint32(n))
}

func produceHash(msgType byte, bytes []byte) common.Hash {
	bytes[0] = msgType
	hashBytes := hasher.Sum(bytes)
	result := common.Hash{}
	result.SetBytes(hashBytes)
	return result
}

var hasher = sha3.NewKeccak256()