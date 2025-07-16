package broadcaster

import (
	servicev1 "github.com/go-pantheon/roma/gen/api/server/broadcaster/service/push/v1"
)

func NewClient(conn *Conn) servicev1.PushServiceClient {
	return servicev1.NewPushServiceClient(conn.ClientConnInterface)
}
