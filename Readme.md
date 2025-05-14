# Project Structure

Below is the structure of the `go-ems` project, reflecting its actual layout:

```
go-ems/
├── cmd/                # Application entry points
├── config/             # Configuration logic
├── conn/               # Database connection setup
├── domain/             # Repository and service interface definitions
├── env/                # Environment-specific configuration files
├── handlers/           # HTTP request handlers (controllers)
├── middlewares/        # Middleware functions for request processing
├── models/             # Data models and structures
├── repositories/       # Database interaction logic
├── routes/             # API route definitions
├── server/             # Server setup and initialization
├── services/           # Core business logic
├── types/              # Shared type definitions
├── utils/              # Shared utility functions
├── config.json         # JSON configuration file
├── Dockerfile          # Docker image definition
├── docker-compose.yml  # Docker Compose configuration
├── Makefile            # Build and automation commands
├── build-n-serve.sh    # Script to build and serve the application
├── go.mod              # Go module definition
├── go.sum              # Dependency checksum file
└── README.md           # Project documentation

# Dependencies

- Consul
- Mysql
- Redis
- Asynqmon

# Run the project

- with config.json
```bash
make run
```

- with env/config.local.json
```bash
make run-local
```

- with docker
```
go mod vendor
docker compose up -d 
```