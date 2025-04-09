package player

import (
	userv1 "github.com/go-pantheon/roma/gen/api/server/player/service/user/v1"
)

func NewUserClient(conn *Conn) userv1.UserServiceClient {
	return userv1.NewUserServiceClient(conn.ClientConnInterface)
}
