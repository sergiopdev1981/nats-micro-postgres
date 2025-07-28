#!/bin/bash

# Script to run go mod tidy and golangci-lint in all microservice directories

# List of directories to process
directories=(
  "add_user"
  "delete_user"
  "get_user"
  "get_users"
  "update_user"
)

# Process each directory
for dir in "${directories[@]}"; do
  echo "Processing $dir..."

  # Enter the directory
  cd "$dir" || { echo "Failed to enter $dir"; exit 1; }

  # Run go mod tidy
  echo "Running go mod tidy in $dir..."
  go mod tidy || { echo "go mod tidy failed in $dir"; cd -; continue; }

  # Run golangci-lint
  echo "Running golangci-lint in $dir..."
  golangci-lint run || { echo "Linting failed in $dir"; cd -; continue; }

  # Return to the root directory
  cd - || { echo "Failed to return to root directory"; exit 1; }

done
