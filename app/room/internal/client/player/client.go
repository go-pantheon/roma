package player

import (
	userv1 "github.com/vulcan-frame/vulcan-game/gen/api/server/player/service/user/v1"
)

func NewUserClient(conn *Conn) userv1.UserServiceClient {
	return userv1.NewUserServiceClient(conn.ClientConnInterface)
}
