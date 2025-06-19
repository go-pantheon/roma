package object

import (
	"maps"

	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
)

const defaultWeightsSize = 8

type Plunder struct {
	Weights map[int64]int64
}

func NewPlunder() *Plunder {
	p := &Plunder{
		Weights: make(map[int64]int64, defaultWeightsSize),
	}
	return p
}

func (o *Plunder) encodeServer(p *dbv1.UserPlunderProto) {
	p.Weights = make(map[int64]int64, len(o.Weights))
	maps.Copy(p.Weights, o.Weights)
}

func (o *Plunder) decodeServer(p *dbv1.UserPlunderProto) {
	if p == nil || p.Weights == nil {
		return
	}

	o.Weights = make(map[int64]int64, len(p.Weights))
	maps.Copy(o.Weights, p.Weights)
}
