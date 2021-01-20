/*
Package ranges implements weighted range distribution algorithms for efficient value selection.
It supports weighted random selection with exclusion capabilities and list generation.
*/

package weights

import (
	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/pkg/util/maths/u64"
	"github.com/vulcan-frame/vulcan-game/pkg/util/ranges"
)

// Weights represents a collection of weighted ranges for probabilistic selection
type Weights struct {
	ranger *ranges.Range
}

func TryNewWeights(weights []uint64, values []int64) (*Weights, error) {
	if len(weights) == 0 {
		return nil, errors.New("new weights from ints failed: ints is empty")
	}
	if len(weights) != len(values) {
		return nil, errors.New("new weights from ints failed: ints and values must have the same length")
	}

	coords := make([]uint64, 0, len(weights)*2)
	start := uint64(0)
	for _, weight := range weights {
		coords = append(coords, start, start+weight)
		start += weight
	}

	rgs, err := ranges.TryNewRange(coords, values)
	if err != nil {
		return nil, err
	}
	return &Weights{
		ranger: rgs,
	}, nil
}

func (ws *Weights) Value(weight uint64) (int64, bool) {
	v, i := ws.ranger.Find(weight)
	return v, i >= 0
}

func (ws *Weights) Rand() int64 {
	if ws.ranger.Len() == 0 {
		return 0
	}

	if ret, i := ws.ranger.Find(u64.Random(uint64(ws.ranger.Max()))); i >= 0 {
		return ret
	}
	v, _ := ws.ranger.Rand()
	return v
}

// RandList generates a random list based on weights, the result is not duplicated.
// It will return less than count if there are not enough values.
func (ws *Weights) RandList(count int) []int64 {
	if count <= 0 || ws.ranger.Len() == 0 {
		return make([]int64, 0)
	}

	result := make(map[int64]struct{}, count)

	max := ws.ranger.Max()
	for len(result) < count {
		if max == 0 {
			break
		}

		w := u64.Random(max)

		var (
			v int64
			i int
		)
		for j := 0; j < ws.ranger.Len(); j++ {
			v, i = ws.ranger.Find(w)
			if i == ranges.NotFound {
				break
			}
			_, exist := result[v]
			if !exist {
				break
			}
			// step to next range
			w += ws.ranger.Pairs[i].Len()
		}
		if i == ranges.NotFound {
			break
		}

		result[v] = struct{}{}
		max -= ws.ranger.Pairs[i].Len()
	}

	ret := make([]int64, 0, len(result))
	for v := range result {
		ret = append(ret, v)
	}
	return ret
}
