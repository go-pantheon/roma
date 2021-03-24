package domain

import (
	"github.com/go-kratos/kratos/log"
)

type PlunderDomain struct {
	log *log.Helper
}

func NewPlunderDomain(logger log.Logger) *PlunderDomain {
	do := &PlunderDomain{
		log: log.NewHelper(log.With(logger, "module", "player/gate/domain/plunder")),
	}
	return do
}
