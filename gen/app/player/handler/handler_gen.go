// Code generated by gen-api. DO NOT EDIT.

package handler

import (
	"context"
	"github.com/go-pantheon/roma/gen/api/client/module"
	"github.com/go-pantheon/roma/gen/app/player/service"
	"github.com/go-pantheon/fabrica-net/xnet"
	"github.com/pkg/errors"
)

func PlayerHandle(ctx context.Context, s *service.PlayerServices, in xnet.TunnelMessage) ([]byte, error) {
	var (
		out []byte
		err error
	)

	switch climod.ModuleID(in.GetMod()) {
	
	case climod.ModuleID_Dev:
		out, err = handleDev(ctx, s, in.GetMod(), in.GetSeq(), in.GetObj(), in.GetData())
	case climod.ModuleID_Hero:
		out, err = handleHero(ctx, s, in.GetMod(), in.GetSeq(), in.GetObj(), in.GetData())
	case climod.ModuleID_Storage:
		out, err = handleStorage(ctx, s, in.GetMod(), in.GetSeq(), in.GetObj(), in.GetData())
	case climod.ModuleID_System:
		out, err = handleSystem(ctx, s, in.GetMod(), in.GetSeq(), in.GetObj(), in.GetData())
	case climod.ModuleID_User:
		out, err = handleUser(ctx, s, in.GetMod(), in.GetSeq(), in.GetObj(), in.GetData())
	default:
		err = errors.Errorf("mod not found. mod=%d", in.GetMod())
	}
	return out, err
}

