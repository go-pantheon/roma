package gamedata

import "github.com/go-pantheon/roma/pkg/util/weights"

type GroupWeights struct {
	*weights.Weights
}

func TryNewGroupWeights(ends []uint64) (*GroupWeights, error) {
	values := make([]int64, len(ends))
	for i := range len(values) {
		values[i] = int64(i)
	}

	ws, err := weights.TryNewWeights(ends, values)
	if err != nil {
		return nil, err
	}

	return &GroupWeights{
		Weights: ws,
	}, nil
}
