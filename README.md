# Graded Challenge 2 - Kevin Sofyan

An application for managing books and borrowing using gRPC and REST APIs. Features user authentication, book management, and automated status updates for borrowed books.

## Prerequisites

- Docker & Docker Compose
- Go 1.21+
- Protocol Buffers compiler
- gRPC tools

## Setup

1. **Clone the repository**:

   ```sh
   git clone https://github.com/kevinsofyan/gc-2-kevinsofyab
   cd gc-2-kevinsofyab
   ```

2. **Install Go dependencies**:

   Ensure you have Go installed and run the following command in each service directory (user-service, product-service, order-service):

   ```sh
   go mod tidy
   ```
## Build and Run the Application

1. **Build Docker images:**
  Use Docker Compose to build the Docker images for all services:
  ```sh
  docker-compose build
   ```

2. **Run Docker containers:**
  Use Docker Compose to run the Docker containers:
  ```sh
  docker-compose up -d
   ```
3. **Check the running containers:**
  Ensure all containers are running:
  ```sh
  docker ps
   ```   