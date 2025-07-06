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

# 🚀 Production Deployment Guide (Docker + EC2 + ECR + Monitoring)

This guide walks through deploying a Dockerized application to an EC2 instance using Docker Compose, pulling images from AWS ECR, and enabling observability with Prometheus and Grafana.

---

## 🛠 Prerequisites

* EC2 instance (Amazon Linux 2 or similar)
* SSH access (`.pem` key file)
* Docker images pushed to AWS ECR
* AWS CLI configured
* Monitoring stack files (`docker-compose.yml`, `prometheus.yml`, etc.)

---

## 📦 Step 1: Setup Docker and Docker Compose

SSH into the EC2 instance:

```bash
ssh -i path/to/your-key.pem ec2-user@your-ec2-hostname
```

Then run:

```bash
sudo dnf update -y
sudo dnf install -y docker
sudo systemctl enable --now docker
sudo usermod -aG docker $USER
newgrp docker
```

Verify Docker is installed:

```bash
docker version
docker info
```

Install Docker Compose plugin:

```bash
mkdir -p ~/.docker/cli-plugins/
curl -SL https://github.com/docker/compose/releases/latest/download/docker-compose-linux-x86_64 \
  -o ~/.docker/cli-plugins/docker-compose
chmod +x ~/.docker/cli-plugins/docker-compose
docker compose version
```

---

## 📂 Step 2: Copy Files to the Server

From your **local machine**, copy required files to the EC2 instance:

```bash
scp -i path/to/your-key.pem ./docker-compose.yml ec2-user@your-ec2-hostname:/home/ec2-user/docker-compose.yml
scp -i path/to/your-key.pem ./init_db.sql ec2-user@your-ec2-hostname:/home/ec2-user/init_db.sql
scp -i path/to/your-key.pem ./prometheus.yml ec2-user@your-ec2-hostname:/home/ec2-user/prometheus.yml
```

---

## 🔐 Step 3: Authenticate Docker with AWS ECR

On the EC2 server:

1. Configure AWS CLI:

```bash
aws configure
```

2. Authenticate Docker to AWS ECR:

```bash
aws ecr get-login-password --region <your-region> \
  | docker login --username AWS --password-stdin <your-account-id>.dkr.ecr.<your-region>.amazonaws.com
```

---

## 🐳 Step 4: Deploy with Docker Compose

Pull latest images:

```bash
docker compose -f docker-compose.yml pull
```

Start the containers:

```bash
docker compose -f docker-compose.yml up -d
```

---

## 📊 Step 5: Access and Configure Grafana

* Visit: `http://<your-ec2-public-ip>:3000`
* Default login: `admin / admin`
* Add **Prometheus** as a data source:

    * URL: `http://prometheus:9090`

### Sample Prometheus Queries

* CPU Usage:

  ```
  sum(rate(process_cpu_seconds_total[1m]))
  ```

* Custom request metric:

  ```
  sum(your_app_requests_total) by (url)
  ```

---

## ✅ Monitoring & Troubleshooting

Check container status:

```bash
docker ps
```

Check logs:

```bash
docker compose logs -f
```

Restart services:

```bash
docker compose restart
```

Shutdown:

```bash
docker compose down
```

---

## 📌 Notes

* Open the necessary ports in your EC2 security group:
    * `22` for SSH
    * `80`/`443` for web access
    * `3000` for Grafana
    * `8080` for Backend service
* Do **not** expose Prometheus publicly without security controls.