package sorts

import (
	"math/rand"
)

func I64Mix(parts []int64) {
	n := len(parts)
	for i := range parts {
		idx := n - i - 1
		swap := rand.Intn(n - i)
		parts[idx], parts[swap] = parts[swap], parts[idx]
	}
}

func U64Mix(parts []uint64) {
	n := len(parts)
	for i := range parts {
		idx := n - i - 1
		swap := rand.Intn(n - i)
		parts[idx], parts[swap] = parts[swap], parts[idx]
	}
}
