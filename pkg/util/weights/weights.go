/*
Package ranges implements weighted range distribution algorithms for efficient value selection.
It supports weighted random selection with exclusion capabilities and list generation.
*/

package weights

import (
	"github.com/go-pantheon/roma/pkg/util/maths/u64"
	"github.com/go-pantheon/roma/pkg/util/ranges"
	"github.com/pkg/errors"
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
	filteredValues := make([]int64, 0, len(values))

	var totalWeight uint64

	for i, weight := range weights {
		if weight > 0 {
			coords = append(coords, totalWeight, totalWeight+weight)
			filteredValues = append(filteredValues, values[i])
		}

		totalWeight += weight
	}

	// if all weights are 0, create an empty ranger
	if len(coords) == 0 {
		return &Weights{ranger: &ranges.Range{}}, nil
	}

	rgs, err := ranges.TryNewRange(coords, filteredValues)
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

	if ret, i := ws.ranger.Find(u64.Random(ws.ranger.Max())); i >= 0 {
		return ret
	}

	v, _ := ws.ranger.Rand()

	return v
}

func (ws *Weights) RandList(count int) []int64 {
	if count <= 0 || ws.ranger.Len() == 0 {
		return make([]int64, 0)
	}

	result := make([]int64, 0, count)
	mx := ws.ranger.Max()

	for range count {
		w := u64.Random(mx)

		if v, i := ws.ranger.Find(w); i != ranges.NotFound {
			result = append(result, v)
		}
	}

	return result
}
