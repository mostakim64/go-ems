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

```

# Tech Stack

- Language : Golang
- Database : Mysql
- Cache: Redis
- Asynq Queue: Asynq
- Config Management: Consul
- Library
  - [Cobra](1) - Framework for building CLI applications
  - [Viper](2) - Library for managing configuration files and environment variables
  - [GORM](3) - ORM library for interacting with relational databases
  - [Echo](4) - Web framework used for building RESTful APIs and handling HTTP routing
  - [Ozzo-Validation](5) - Library used for validating input data, ensuring data integrity and consistency
  - [golang-course-utils](6) - Utility library providing shared helper functions and reusable components
  - [Asynq](7) - Library for managing background tasks and distributed task queues

# Using Private Go Packages

To use private Go packages in your project, follow these steps:

1. **Configure Git to use SSH for private repositories**:
   Add the following to your `~/.gitconfig` file:
   ```bash
   [url "ssh://git@github.com/"]
       insteadOf = https://github.com/
   ```

2. **Set the `GOPRIVATE` environment variable**:
   Add the private repository domain to the `GOPRIVATE` environment variable. For example:
   ```bash
   export GOPRIVATE=github.com/your-org-name
   ```

3. **Authenticate with SSH**:
   Ensure your SSH key is added to your SSH agent and linked to your GitHub account:
   ```bash
   ssh-add ~/.ssh/id_rsa
   ```

4. **Get the private package**:
   Use `go get` to fetch the private package:
   ```bash
   go get github.com/your-org-name/private-repo-name
   ```

This setup ensures that Go modules can fetch private repositories securely using SSH.

# Install dependencies

```bash
go mod tidy
go mod vendor
```

# Run the project

> Make sure consul is up and running

## Publish config to consul

```bash
curl --request PUT \
    --data-binary @"$CONFIG_PATH" \
    "$CONSUL_URL/v1/kv/$CONSUL_PATH"
```

## Export environment variables

```bash
export CONSUL_URL="$CONSUL_URL"
export CONSUL_PATH="$CONSUL_PATH"
```

## Build & Run
- build the project
```bash
go build -o app .
```

## Run server

```bash
./app serve
```

## Makefile
- with config.json
```bash
make run
```

- with env/config.local.json
```bash
make run-local
```

## Docker
```
go mod vendor
docker compose up -d 
```

[7]: https://github.com/hibiken/asynq
[6]: https://github.com/vivasoft-ltd/golang-course-utils
[5]: https://github.com/go-ozzo/ozzo-validation
[4]: https://echo.labstack.com/docs/quick-start
[3]: https://gorm.io/docs/index.html
[2]: https://github.com/spf13/viper
[1]: https://github.com/spf13/cobra