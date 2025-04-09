# roma-exceldata

**roma-exceldata** is an Excel configuration file example for the game service [roma](https://github.com/go-pantheon/roma), used to demonstrate game data configuration methods.

## Overview

This repository serves as a reference implementation for managing game configurations using Excel files in the Roma framework. It demonstrates best practices for structuring game data and showcases the powerful features of Roma's configuration system.

## Features

- Supports 2D tables and K-V tables
- Supports multiple data formats, such as int, float, string, bool, array, map, etc.
- Supports data table relationships, including one-to-many, many-to-one, and many-to-many
- Supports data table inheritance
- Supports various merging methods for specified fields
- Supports configuration data validation rules
- Supports custom field generation types

## Installation

Clone the repository:
```bash
git clone https://github.com/go-pantheon/roma-exceldata.git
cd roma-exceldata
```

## Usage

### Data Generation
Uses [vulcan](https://github.com/go-pantheon/roma/tree/main/vulcan/app/gamedata/) in [roma](https://github.com/go-pantheon/roma) project to convert Excel files to JSON files and generate Go code:

```bash
# must in root directory of roma project
make tools

# generate json and go code
make gen-all-data

# only generate json (option) 
make gen-data-json

# only generate go code (option) 
make gen-data-base && make gen-datas
```

### Example Configuration

The `examples` directory contains sample Excel files demonstrating various configuration patterns:

- `Hero/Hero.xlsx`: the configuration of collectable heroes in the game
- `Property/Attribute.xlsx`: the configuration of attributes and calculation formulas in the game
- `Recharge/Recharge.xlsx`: the configuration of platforms and recharge in the game
- `Resource/Resource.xlsx`: the configuration of items and gift packs in the game

## Project Structure

```
.
├── bin/            # Game tool executable file of configuration generation and check
├── excel/          # Excel configuration files
└── server/         # Game server executable file for game designers
```

## go-pantheon

**go-pantheon** is an out-of-the-box game server framework providing a high-performance, highly available game server cluster solution based on microservices architecture.

### Core Features

- 🚀 Microservices game server architecture built with [go-kratos](https://github.com/go-kratos/kratos)
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

## Contributing

We welcome contributions! Please submit any suggestions via [issues](https://github.com/go-pantheon/janus/issues) or [pull requests](https://github.com/go-pantheon/janus/pulls).

## Version

Current version: v0.1.0

## License

This project is licensed under the [MIT License](LICENSE).
