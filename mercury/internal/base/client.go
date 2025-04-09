package base

import (
	"github.com/go-pantheon/roma/gen/api/server/player/admin/user/v1"
	gg "google.golang.org/grpc"
)

type AdminClients struct {
	UserCli adminv1.UserAdminClient
}

func NewAdminClients(conn gg.ClientConnInterface) *AdminClients {
	clis := &AdminClients{
		UserCli: adminv1.NewUserAdminClient(conn),
	}

	return clis
}
