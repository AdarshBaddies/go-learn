# Go-Learn GraphQL Microservices Project

![Go Logo](https://golang.org/lib/godoc/images/go-logo-blue.svg)

## Overview

This project demonstrates a simple, modular microservices architecture built in Go, integrated via gRPC, and exposed through a GraphQL API. It consists of three core services:

1. **Account Service** - Manages user accounts.
2. **Catalog Service** - Handles product catalog data.
3. **Order Service** - Processes orders and tracks ordered products.

The GraphQL layer aggregates these services to provide a unified API endpoint.

## Technologies Used

| Technology | Logo                                                                                  |
| ---------- | ------------------------------------------------------------------------------------- |
| Go         | ![Go Logo](https://golang.org/lib/godoc/images/go-logo-blue.svg)                      |
| GraphQL    | ![GraphQL Logo](https://upload.wikimedia.org/wikipedia/commons/1/17/GraphQL_Logo.svg) |
| gRPC       | ![gRPC Logo](https://grpc.io/img/grpc-logo.png)                                       |
| Docker     | ![Docker Logo](https://www.docker.com/sites/default/files/d8/2019-07/Moby-logo.png)   |
| PostgreSQL | ![Postgres Logo](https://www.postgresql.org/media/img/about/press/elephant.png)       |

## Project Structure

```
account/      # Account service (gRPC + DB)
catalog/      # Catalog service (gRPC)
graphql/      # GraphQL gateway and resolvers
order/        # Order service (gRPC + DB)
docker-compose.yaml  # For running all services together
```

## Features

* Create and fetch accounts
* Manage products
* Create orders with multiple products
* GraphQL API for querying and mutation
* Pagination support

## Prerequisites

* Docker & Docker Compose
* Go 1.21+ (for local development if needed)

## How to Run

1. Clone the repository:

```bash
git clone <repo-url>
cd go-learn
```

2. Start all services using Docker Compose:

```bash
docker-compose up --build
```

3. Access GraphQL Playground in your browser:

```
http://localhost:8080/playground
```

4. Send queries and mutations as defined in `graphql/schema.graphql`.

## Example Queries

```graphql
query {
  accounts {
    id
    name
    orders {
      id
      totalPrice
    }
  }
}
```

```graphql
mutation {
  createAccount(account: { name: "John Doe" }) {
    id
    name
  }
}
```

## How It Was Built

1. Defined Protobuf files for each service.
2. Generated Go gRPC clients and servers.
3. Implemented service logic (DB interaction and business logic).
4. Built GraphQL layer with `gqlgen` to combine services.
5. Used Docker and Docker Compose to run services together.

## Notes

* Each service runs independently on its own port.
* PostgreSQL databases are initialized via Docker volumes.
* GraphQL server communicates with services via gRPC clients.

---


