# Mixturka — Recipe Microservice

Mixturka is a Go microservice for receiving, storing, and retrieving recipes. It is a final project and is currently under development.

## Overview

The service is designed to:

- Consume recipe messages from a Kafka topic.
- Store recipes and ingredients in PostgreSQL.
- Return a recipe and its ingredients by ID through a gRPC API.
- Provide service version information through an HTTP endpoint.

## Planned Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **Message broker:** Apache Kafka
- **API:** gRPC and HTTP
- **HTTP framework:** Gin
- **Logging:** Zap
- **Database migrations:** Goose
- **Containerization:** Docker and Docker Compose

## Architecture

The project follows Clean Architecture principles and is organized into three main layers:

- **Domain:** business entities.
- **Application:** business logic and use cases.
- **Infrastructure:** database access, Kafka consumer, and API servers.

## Status

The project is in the early stages of development. Features and setup instructions will be added as implementation progresses.
