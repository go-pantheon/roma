package object

import (
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
)

type Plunders struct {
	States map[int64]*PlunderState
}

func NewPlunders() *Plunders {
	p := &Plunders{
		States: make(map[int64]*PlunderState),
	}
	return p
}

func NewPlundersProto() *dbv1.PlundersProto {
	p := &dbv1.PlundersProto{}
	p.Plunders = make(map[int64]*dbv1.PlunderStateProto, 8)
	return p
}

func (os *Plunders) EncodeServer(p *dbv1.PlundersProto) {
	p.Plunders = make(map[int64]*dbv1.PlunderStateProto, len(os.States))
	for k, o := range os.States {
		p.Plunders[k] = NewPlunderStateProto()
		o.EncodeServer(p.Plunders[k])
	}
}

type PlunderState struct {
	Weights map[int64]int64
}

func NewPlunderState() *PlunderState {
	p := &PlunderState{
		Weights: make(map[int64]int64, 8),
	}
	return p
}

func NewPlunderStateProto() *dbv1.PlunderStateProto {
	p := &dbv1.PlunderStateProto{}
	return p
}

func (o *PlunderState) EncodeServer(p *dbv1.PlunderStateProto) {
	p.Weights = o.Weights
}

func (os *Plunders) DecodeServer(ps *dbv1.PlundersProto) {
	if ps == nil {
		return
	}

	for k, p := range ps.Plunders {
		o := NewPlunderState()
		o.DecodeServer(p)
		os.States[k] = o
	}
}

func (o *PlunderState) DecodeServer(p *dbv1.PlunderStateProto) {
	if p == nil || p.Weights == nil {
		return
	}

	o.Weights = p.Weights
}
