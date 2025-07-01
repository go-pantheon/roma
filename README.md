# Roma

Roma is a high-performance game server framework based on microservice architecture, developed in Go. This framework aims to provide scalable game server infrastructure that supports various game types. Roma is a core component of the go-pantheon ecosystem, responsible for implementing specific game logic.

For more, please check out: [deepwiki/go-pantheon/roma](https://deepwiki.com/go-pantheon/roma)

## go-pantheon Ecosystem

**go-pantheon** is an out-of-the-box game server framework providing a high-performance, highly available game server cluster solution based on microservices architecture. Roma, as the game logic implementation component, works alongside other core services to form a complete game service ecosystem:

- **Roma**: Game core logic services
- **Janus**: Gateway service for client connection handling and request forwarding
- **Lares**: Account service for user authentication and account management
- **Senate**: Backend management service providing operational interfaces

### Core Features

- 🚀 Microservice game server architecture built with [go-kratos](https://github.com/go-kratos/kratos)
- 🔒 Multi-protocol support (TCP/KCP/WebSocket)
- 🛡️ Enterprise-grade secure communication protocol
- 📊 Real-time monitoring & distributed tracing
- 🔄 Gray release & hybrid deployment support
- 🔍 Developer-friendly debugging environment

### Service Layer Features

- gRPC for inter-service communication
- Stateful dynamic routing & load balancing
- Canary release & hybrid deployment support
- Hot updates & horizontal scaling without downtime

## System Architecture

The relationship between Roma and other go-pantheon components is illustrated below:

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │────▶│  Janus GW   │────▶│  Lares Acct │
└─────────────┘     └──────┬──────┘     └─────────────┘
                           │                    ▲
                           ▼                    │
                    ┌─────────────┐     ┌─────────────┐
                    │  Roma Game  │────▶│Senate Admin │
                    └─────────────┘     └─────────────┘
```

Roma internally adopts a microservice architecture, with services communicating via gRPC:

```
                    ┌─────────────────────────────┐
                    │           Roma              │
                    │                             │
┌─────────────┐     │  ┌───────┐     ┌────────┐  │
│  Config     │◀───▶│  │Player │◀───▶│ Room   │  │
│  (etcd)     │     │  │Service│     │Service │  │
└─────────────┘     │  └───────┘     └────────┘  │
                    │       ▲             ▲      │
┌─────────────┐     │       │             │      │
│  Monitoring │◀───▶│       └─────┬───────┘      │
│(Prometheus) │     │             │              │
└─────────────┘     │        ┌────▼────┐         │
                    │        │ Shared  │         │
┌─────────────┐     │        │ Cache   │         │
│  Tracing    │◀───▶│        │(Redis)  │         │
│     (OT)    │     │                             │
└─────────────┘     └─────────────────────────────┘
```

## Project Overview

The Roma framework is built on the Go-Kratos microservice framework, supporting both gRPC and HTTP protocols, and integrates core components such as etcd registry, Redis cache, and MongoDB database. The framework design follows Domain-Driven Design (DDD) principles, achieving high cohesion and low coupling of business logic through clear service boundaries and domain models.

## Technology Stack

Roma utilizes the following core technologies:

| Technology/Component | Purpose                      | Version |
| -------------------- | ---------------------------- | ------- |
| Go                   | Primary development language | 1.24+   |
| go-kratos            | Microservice framework       | v2.8.4  |
| gRPC                 | Inter-service communication  | v1.71.1 |
| Protobuf             | Data serialization           | v1.36.6 |
| etcd                 | Service discovery & registry | v3.5.21 |
| Redis                | Caching                      | v9.7.3  |
| MongoDB              | Data storage                 | v2.1.0  |
| OpenTelemetry        | Distributed tracing          | v1.35.0 |
| Prometheus           | Monitoring system            | v1.22.0 |
| Google Wire          | Dependency injection         | v0.6.0  |
| Buf                  | API management               | Latest  |

## Key Features

- **Microservice Architecture**: Built on Go-Kratos with service registry, discovery, and load balancing
- **Multi-Protocol Support**: Simultaneous support for gRPC and HTTP interfaces
- **Configuration Center**: etcd-based configuration center with dynamic updates
- **Game Data Configuration**: Support for Excel-format game data import and code generation
- **Distributed Tracing**: OpenTelemetry integration for distributed tracing
- **Service Monitoring**: Prometheus metrics collection
- **Dependency Injection**: Google Wire for dependency injection
- **Code Generation**: Simplified development through Protobuf and code generation tools

## Core Components

### Application Services (app/)

- **player**: Business functionality related to player progression
- **room**: Business functionality related to battle rooms

### API Definitions (api/)

- **client**: Client API definitions
- **server**: Server-side internal API definitions
- **db**: Database model definitions

### Testing Tools (mercury/)

Client connection simulation tool for server functionality testing and performance testing

### Toolchain (vulcan/)

Code generation tool based on api/ and exceldata/ directories, generating API connector code and gamedata structures with basic logic

### Game Data (gamedata/)

Game configuration data and processing logic

### Common Libraries (pkg/)

- **util**: General utility functions
- **errs**: Error handling
- **universe**: Common business logic

## Requirements

- Go 1.24+
- Protobuf
- etcd
- Redis
- MongoDB

## Quick Start

### Initialize Environment

```bash
make init
```

### Generate API Code

```bash
make proto
make api
```

### Generate Game Data

```bash
make gen-all-data
```

### Build Services

```bash
make build
```

### Start Services

```bash
# Start all services
make run

# Start a specific service
make run app=player
```

## Integration with go-pantheon Components

Integration of Roma with other go-pantheon components typically follows these steps:

### Integration with Janus Gateway

1. Configure Roma service registry information to ensure discovery by Janus
2. Set up Roma service routing rules in Janus
3. Configure load balancing strategies

```yaml
# Janus configuration example
services:
  - name: player
    discovery:
      type: etcd
      address: ["127.0.0.1:2379"]
    endpoints:
      - protocol: grpc
        port: 9000
      - protocol: http
        port: 8000
```

### Integration with Lares Account System

1. Client accesses the Lares account service to obtain an AuthToken
2. When establishing a connection with the Janus gateway, the client sends the AuthToken
3. After validating the AuthToken, the Janus gateway establishes a connection with the client and creates a Tunnel with Roma services, forwarding subsequent messages

```proto
# api/
message AuthToken {
  int32 location = 1; // Location
  int64 account_id = 2; // Account ID
  string rand = 3; // random string for replay attack prevention
  int64 timeout = 4; // timeout (unix timestamp)
  bool unencrypted = 5; // Whether the token is unencrypted, default is false
  string color = 6; // Color for route distribution
  OnlineStatus status = 7; // Online status
}
```

### Integration with Senate Backend Management

1. Ensure Roma services expose necessary management APIs
2. Call `api/server/**/admin` interfaces in the Senate service

## Project Structure

```
.
├── api/               # API definitions
│   ├── client/        # Client API
│   ├── db/            # Database models
│   └── server/        # Server API
├── app/               # Application services
│   ├── player/        # Player service
│   └── room/          # Room service
├── bin/               # Compiled output directory
├── buf/               # Buf configuration files
├── deps/              # Local dependencies
├── exceldata/         # Excel game configuration data
├── gamedata/          # Game data processing
├── gen/               # Generated code
├── mercury/           # Client simulation testing tool
├── pkg/               # Common libraries
│   ├── errs/          # Error handling
│   ├── universe/      # Common business logic
│   └── util/          # Utility functions
├── third_party/       # Third-party dependencies
└── vulcan/            # API and game data code generation tools
```

## Port Conventions

### Player and Room Services

- TCP Ports:
  - External: 70xx
- HTTP Ports:
  - Internal: 81xx
  - External: 80xx
- gRPC Ports:
  - Internal: 91xx
  - External: 90xx

## Development Guide

### Development Workflow

1. Define API interfaces using Protobuf
2. Generate interface code with `make proto` and `make api`
3. Implement service logic based on business requirements
4. Use Wire for dependency injection
5. Write unit tests
6. Build and deploy services

### Adding New Services

Steps to create a new service:

1. Create Proto definitions for the new service in the `api/server/` directory
2. Generate API code: `make proto`
3. Create a new service directory in `app/`
4. Copy and modify the framework code from existing services
5. Generate dependency injection code using Wire: `make wire`
6. Implement service logic

### Game Data Configuration

1. Create Excel configuration files in the `exceldata/` directory
2. Generate game data code: `make gen-all-data`
3. Use generated data structures in services

## Troubleshooting

### 1. Service Registration Failure

**Issue**: Service cannot register with etcd

**Solution**:
- Check if etcd is running properly
- Verify that the etcd address in the configuration file is correct
- Check network connectivity

### 2. Code Generation Errors

**Issue**: Generated code has errors after running `make proto`

**Solution**:
- Ensure all necessary protoc plugins are installed
- Check proto file syntax
- Run `make init` to reinstall all dependency tools

### 3. Configuration Errors During Service Startup

**Issue**: Service fails to start with configuration errors

**Solution**:
- Check configuration files in the `configs/` directory
- Ensure all necessary environment variables are set
- Reference template files in the `configs.tmpl/` directory

## Contributing

1. Fork this repository
2. Create a feature branch
3. Submit changes
4. Ensure all tests pass
5. Submit a Pull Request

## License

This project is licensed under the terms specified in the LICENSE file.
