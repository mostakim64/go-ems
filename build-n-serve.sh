#!/bin/bash

# filepath: /Users/shuvojitsaha/Desktop/vivasoft/golang-course/go-ems/build_and_run.sh

# Exit immediately if a command exits with a non-zero status
set -e

# Default values
CONSUL_URL="http://127.0.0.1:8500"
CONSUL_PATH="event-management"
CONFIG_PATH="config.json"

# Function to display usage
usage() {
    echo "Usage: $0 [-c <config_path>]"
    echo "  -c <config_path>   Path to the JSON configuration file (default: env/config.local.json)"
    exit 1
}

# Parse command-line arguments
while getopts "c:" opt; do
    case $opt in
        c) CONFIG_PATH="$OPTARG" ;;
        *) usage ;;
    esac
done

# Step 1: Push configuration to Consul
echo "Pushing configuration to Consul..."
if [[ ! -f "$CONFIG_PATH" ]]; then
    echo "Error: Configuration file not found at $CONFIG_PATH"
    exit 1
fi

curl --request PUT \
    --data-binary @"$CONFIG_PATH" \
    "$CONSUL_URL/v1/kv/$CONSUL_PATH"

echo "Configuration pushed to Consul at $CONSUL_URL/v1/kv/$CONSUL_PATH"

# Step: 2: Install dependencies
echo "Installing dependencies..."
go mod tidy
go mod vendor

# Step 3: Build the Go project
echo "Building the Go project..."
go build -o app .

echo "Build completed successfully."

# Step 4: Set environment variables
export CONSUL_URL="$CONSUL_URL"
export CONSUL_PATH="$CONSUL_PATH"

echo "Environment variables set:"
echo "  CONSUL_URL=$CONSUL_URL"
echo "  CONSUL_PATH=$CONSUL_PATH"

# Step 5: Run the project
echo "Running the project..."
./app serve