package life

import (
	"context"

	intrav1 "github.com/vulcan-frame/vulcan-game/gen/api/server/gate/intra/v1"
	"github.com/vulcan-frame/vulcan-kit/xcontext"
)

func IsInnerContext(ctx context.Context) bool {
	return IsInnerStatus(OnlineStatus(xcontext.Status(ctx)))
}

func IsInnerStatus(status intrav1.OnlineStatus) bool {
	return IsSvcStatus(status) || IsAdminStatus(status)
}

func IsDevContext(ctx context.Context) bool {
	return IsDevStatus(OnlineStatus(xcontext.Status(ctx)))
}

func IsSvcContext(ctx context.Context) bool {
	return IsSvcStatus(OnlineStatus(xcontext.Status(ctx)))
}

func IsAdminContext(ctx context.Context) bool {
	return IsAdminStatus(OnlineStatus(xcontext.Status(ctx)))
}

func IsGateContext(ctx context.Context) bool {
	return IsGateStatus(OnlineStatus(xcontext.Status(ctx)))
}

func IsSvcStatus(status intrav1.OnlineStatus) bool {
	return status == intrav1.OnlineStatus_ONLINE_STATUS_SVC
}

func IsAdminStatus(status intrav1.OnlineStatus) bool {
	return status == intrav1.OnlineStatus_ONLINE_STATUS_ADMIN
}

func IsGateStatus(status intrav1.OnlineStatus) bool {
	return status == intrav1.OnlineStatus_ONLINE_STATUS_GATE
}

func IsDevStatus(status intrav1.OnlineStatus) bool {
	return status == intrav1.OnlineStatus_ONLINE_STATUS_DEV
}

func OnlineStatus(status int64) intrav1.OnlineStatus {
	return intrav1.OnlineStatus(status)
}
