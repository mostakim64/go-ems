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
* Do **not** expose Prometheus publicly without security controls.