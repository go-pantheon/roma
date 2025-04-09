// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-net/health"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds/cmdregistrar"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds/storage"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds/system"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds/user"
	registry5 "github.com/go-pantheon/roma/app/player/internal/app/dev/gate/registry"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/service"
	biz7 "github.com/go-pantheon/roma/app/player/internal/app/gamedata/admin/biz"
	registry9 "github.com/go-pantheon/roma/app/player/internal/app/gamedata/admin/registry"
	service9 "github.com/go-pantheon/roma/app/player/internal/app/gamedata/admin/service"
	biz2 "github.com/go-pantheon/roma/app/player/internal/app/hero/gate/biz"
	domain3 "github.com/go-pantheon/roma/app/player/internal/app/hero/gate/domain"
	registry7 "github.com/go-pantheon/roma/app/player/internal/app/hero/gate/registry"
	service2 "github.com/go-pantheon/roma/app/player/internal/app/hero/gate/service"
	domain4 "github.com/go-pantheon/roma/app/player/internal/app/plunder/gate/domain"
	biz9 "github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/biz"
	data5 "github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/data"
	domain6 "github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/domain"
	registry11 "github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/registry"
	service11 "github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/service"
	biz8 "github.com/go-pantheon/roma/app/player/internal/app/storage/admin/biz"
	registry10 "github.com/go-pantheon/roma/app/player/internal/app/storage/admin/registry"
	service10 "github.com/go-pantheon/roma/app/player/internal/app/storage/admin/service"
	biz3 "github.com/go-pantheon/roma/app/player/internal/app/storage/gate/biz"
	domain2 "github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain"
	registry6 "github.com/go-pantheon/roma/app/player/internal/app/storage/gate/registry"
	service3 "github.com/go-pantheon/roma/app/player/internal/app/storage/gate/service"
	biz4 "github.com/go-pantheon/roma/app/player/internal/app/system/gate/biz"
	registry3 "github.com/go-pantheon/roma/app/player/internal/app/system/gate/registry"
	service4 "github.com/go-pantheon/roma/app/player/internal/app/system/gate/service"
	biz6 "github.com/go-pantheon/roma/app/player/internal/app/user/admin/biz"
	data4 "github.com/go-pantheon/roma/app/player/internal/app/user/admin/data"
	domain5 "github.com/go-pantheon/roma/app/player/internal/app/user/admin/domain"
	registry8 "github.com/go-pantheon/roma/app/player/internal/app/user/admin/registry"
	service8 "github.com/go-pantheon/roma/app/player/internal/app/user/admin/service"
	biz5 "github.com/go-pantheon/roma/app/player/internal/app/user/gate/biz"
	data2 "github.com/go-pantheon/roma/app/player/internal/app/user/gate/data"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain"
	registry4 "github.com/go-pantheon/roma/app/player/internal/app/user/gate/registry"
	service5 "github.com/go-pantheon/roma/app/player/internal/app/user/gate/service"
	"github.com/go-pantheon/roma/app/player/internal/client"
	"github.com/go-pantheon/roma/app/player/internal/client/gate"
	"github.com/go-pantheon/roma/app/player/internal/conf"
	"github.com/go-pantheon/roma/app/player/internal/core"
	"github.com/go-pantheon/roma/app/player/internal/data"
	"github.com/go-pantheon/roma/app/player/internal/intra/filter"
	registry2 "github.com/go-pantheon/roma/app/player/internal/intra/registry"
	service7 "github.com/go-pantheon/roma/app/player/internal/intra/service"
	"github.com/go-pantheon/roma/app/player/internal/server"
	"github.com/go-pantheon/roma/app/player/internal/server/registry"
	service6 "github.com/go-pantheon/roma/gen/app/player/service"
	data3 "github.com/go-pantheon/roma/pkg/universe/data"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(confServer *conf.Server, label *conf.Label, recharge *conf.Recharge, confRegistry *conf.Registry, confData *conf.Data, logger log.Logger, healthServer *health.Server) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	userRepo, err := data2.NewUserMongoRepo(dataData, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	userProtoCache := data2.NewProtoCache()
	userDomain := domain.NewUserDomain(userRepo, logger, userProtoCache)
	routeTable := gate.NewRouteTable(dataData)
	discovery, err := client.NewDiscovery(confRegistry)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	conn, err := gate.NewConn(logger, routeTable, discovery)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	pushServiceClient := gate.NewClient(conn)
	v, err := gate.NewConns(logger, routeTable, discovery)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	v2 := gate.NewClients(v)
	pushRepo := data3.NewPushRepo(pushServiceClient, v2, logger)
	manager, cleanup2 := core.NewManager(logger, userDomain, pushRepo)
	httpFilter := filter.NewHttpFilter(manager)
	servicelessUseCase := registry.NewServicelessUseCase()
	devUseCase := biz.NewDevUseCase(manager, logger)
	showTimeCommander := system.NewShowTimeCommander(devUseCase)
	changeTimeCommander := system.NewChangeTimeCommander(devUseCase)
	storageDomain := domain2.NewStorageDomain(logger)
	addItemCommander := storage.NewAddItemCommander(devUseCase, storageDomain)
	subItemCommander := storage.NewSubItemCommander(devUseCase, storageDomain)
	addPackCommander := storage.NewAddPackCommander(devUseCase, storageDomain)
	subPackCommander := storage.NewSubPackCommander(devUseCase, storageDomain)
	simulateRechargeCommander := user.NewSimulateRechargeCommander(devUseCase)
	createAdminPlayerCommander := user.NewAdminPlayerCommander(devUseCase, storageDomain)
	registrar := cmdregistrar.NewRegistrar(showTimeCommander, changeTimeCommander, addItemCommander, subItemCommander, addPackCommander, subPackCommander, simulateRechargeCommander, createAdminPlayerCommander)
	devServiceServer := service.NewDevService(logger, devUseCase, registrar)
	heroDomain := domain3.NewHeroDomain(logger)
	heroUseCase := biz2.NewHeroUseCase(manager, logger, heroDomain, storageDomain)
	heroServiceServer := service2.NewHeroService(logger, heroUseCase)
	plunderDomain := domain4.NewPlunderDomain(logger)
	storageUseCase := biz3.NewStorageUseCase(manager, logger, storageDomain, plunderDomain)
	storageServiceServer := service3.NewStorageService(logger, storageUseCase)
	systemUseCase := biz4.NewSystemUseCase(manager, logger)
	systemServiceServer := service4.NewSystemService(logger, systemUseCase)
	userUseCase := biz5.NewUserUseCase(manager, userDomain, storageDomain, logger)
	userServiceServer := service5.NewUserService(logger, userUseCase)
	playerServices := service6.NewPlayerServices(devServiceServer, heroServiceServer, storageServiceServer, systemServiceServer, userServiceServer)
	tunnelServiceServer := service7.NewTunnelService(logger, manager, playerServices)
	intraRegistrar := registry2.NewIntraRegistrar(tunnelServiceServer)
	serviceRegistrars := registry.NewServiceRegistrars(servicelessUseCase, intraRegistrar)
	systemRegistrar := registry3.NewSystemRegistrar(systemServiceServer)
	userRegistrar := registry4.NewUserRegistrar(userServiceServer)
	devRegistrar := registry5.NewDevRegistrar(devServiceServer)
	storageRegistrar := registry6.NewStorageRegistrar(storageServiceServer)
	heroRegistrar := registry7.NewHeroRegistrar(heroServiceServer)
	gateRegistrars := registry.NewGateRegistrars(systemRegistrar, userRegistrar, devRegistrar, storageRegistrar, heroRegistrar)
	domainUserRepo, err := data4.NewUserMongoRepo(dataData, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	domainUserDomain := domain5.NewUserDomain(domainUserRepo, logger)
	bizUserUseCase := biz6.NewUserUseCase(manager, domainUserDomain, logger)
	userAdminServer := service8.NewUserAdmin(logger, bizUserUseCase)
	registryUserRegistrar := registry8.NewUserRegistrar(userAdminServer)
	gamedataUseCase := biz7.NewGamedataUseCase(manager, logger)
	gamedataAdmin := service9.NewGamedataAdmin(logger, gamedataUseCase)
	gamedataRegistrar := registry9.NewGamedataRegistrar(gamedataAdmin)
	bizStorageUseCase := biz8.NewStorageUseCase(logger, storageDomain)
	storageAdminServer := service10.NewStorageAdmin(logger, manager, bizStorageUseCase)
	registryStorageRegistrar := registry10.NewStorageRegistrar(storageAdminServer)
	orderRepo, err := data5.NewOrderMongoRepo(dataData, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	orderDomain := domain6.NewOrderDomain(orderRepo, logger)
	rechargeUseCase := biz9.NewRechargeUseCase(manager, orderDomain, logger)
	rechargeAdminServer := service11.NewRechargeAdmin(logger, rechargeUseCase)
	rechargeRegistrar := registry11.NewRechargeRegistrar(rechargeAdminServer)
	adminRegistrars := registry.NewAdminRegistrars(registryUserRegistrar, gamedataRegistrar, registryStorageRegistrar, rechargeRegistrar)
	httpServer := server.NewHTTPServer(confServer, logger, httpFilter, serviceRegistrars, gateRegistrars, adminRegistrars)
	grpcFilter := filter.NewGrpcFilter(manager)
	grpcServer := server.NewGRPCServer(confServer, logger, grpcFilter, serviceRegistrars, adminRegistrars)
	registryRegistrar, err := server.NewRegistrar(confRegistry)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	app := newApp(logger, httpServer, grpcServer, healthServer, label, registryRegistrar)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}
