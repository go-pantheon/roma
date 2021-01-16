module github.com/vulcan-frame/vulcan-game

go 1.15

require (
	github.com/google/wire v0.4.0
	github.com/pkg/errors v0.9.1
	github.com/redis/go-redis v6.15.9
	github.com/stretchr/testify v1.7.0
	github.com/vulcan-frame/vulcan-pkg-app v0.0.0
	github.com/vulcan-frame/vulcan-pkg-tool v0.0.0
	go.uber.org/atomic v1.7.0
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/vulcan-frame/vulcan-pkg-tool => ./pkg/vulcan-pkg-tool
replace github.com/vulcan-frame/vulcan-pkg-app => ./pkg/vulcan-pkg-app
