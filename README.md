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

 ## Endpoints

 ### Register User
```sh
curl -X POST http://localhost:8081/register \
-H "Content-Type: application/json" \
-d '{
    "username": "testuser",
    "password": "testpass123"
}'
```

### Login
```sh
curl -X POSThttp://localhost:8081/login \
-H "Content-Type: application/json" \
-d '{
    "username": "testuser",
    "password": "testpass123"
}'
```  

### Books (protected)

#### Create Book
```sh
curl -X POST http://localhost:8081/books \
-H "Content-Type: application/json" \
-H "Authorization: Bearer {token}" \
-d '{
    "title": "The Go Programming Language",
    "author": "bebas",
    "published_date": "2015-11-05T00:00:00Z"
}'
```

#### Update Book
```sh
curl -X PUT http://localhost:8081/books/65f2e1234567890abcdef124 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer {token}" \
-d '{
    "title": "The Go Programming Language",
    "author": "bebaslah",
    "published_date": "2015-11-05T00:00:00Z"
}'
```

#### Delete Book
```sh
curl -X DELETE http://localhost:8081/books/65f2e1234567890abcdef124 \
-H "Authorization: Bearer {token}"
```

### Borrowed Books (protected)

#### Borrow Book
```sh
curl -X POST http://localhost:8081/borrowed-books/borrow/65f2e1234567890abcdef124 \
-H "Authorization: Bearer {token}"
```

#### Return Book
```sh
curl -X POST http://localhost:8081/borrowed-books/return/65f2e1234567890abcdef125 \
-H "Authorization: Bearer {token}"
```
