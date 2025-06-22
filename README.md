[![Go Report Card](https://goreportcard.com/badge/github.com/craftaholic/insider)](https://goreportcard.com/report/github.com/craftaholic/insider)
[![GoDoc](https://godoc.org/github.com/craftaholic/insider?status.svg)](https://godoc.org/github.com/craftaholic/insider)
[![License](https://img.shields.io/github/license/craftaholic/insider)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/craftaholic/insider)](go.mod)

# Overview
This project is an interview test from Insider for Senior Golang Developer role. It implements an automated message sending system with REST API endpoints for control and monitoring.

# Analyzing
For start and stop feature. I will run the automated service in a dedicated thread. When start/stop signal sent
it will cancel or create a new go routine for that specific automation service (decoulple the lifecycle of the
rest api request with the service running in background)

The requirements only says process 2 message every 2 mins. But I have extended so this project can process way more messages than required.
- Every 2 minutes (Using tick lib from the parent routine - This always run when the service started) -> Push new messages into a channel that has buffer.
- There will be a set of workers - each is a dedicated go routine that listen to the channel above. This will make all message handled concurently.
- To avoid getting the same message from the db, I will create a postgres function that will get oldest pending messages and lock those rows and change the status from pending -> proccessing. This will avoid messages being handled multiple times.

>Note: This design is to get at-least once pattern. If we need exactly once -> should use event-driven.

There are still small edge-cases where by the notification sent but not updated to DB (network/sudden death issue). To tackle this problem, I can do some prediodically scanning for messages in failed/processing status for longer than X amount of minutes and handled them there. But that is out of scope for this project.

# Architecture Design

Overview Logical/Sequence diagram:
![Sequence Diagram](./docs/diagram.png)

Database design:
![Database Design](./docs/database.png)

# Code Structure
```
.
├── .golangci.yaml                      # Store linting config
├── .air.toml                           # Hot reload config
├── .env.example                        # Sample env
├── Dockerfile                          # Application Dockerfile
├── LICENSE
├── README.md
├── build                               # This is the DB initialize file content
│   └── init.sql
├── cmd
│   └── server
│       └── main.go                     # Main.go file - entrypoint of the server
├── devbox.json                         # Development env configuration file - similar to package.json
├── devbox.lock
├── docker-compose.yml                  # Docker-compose for the whole system
├── docker-compose.dev.yml              # Docker-compose for dev env (without the application container)
├── docs                                # Docs related (design docs, swagger etc)
├── e2e                                 # e2e testing
├── internal                            # Main logic
    ├── api                             # API management
    │   ├── middleware
    │   │   └── logging.go
    │   └── route
    │       ├── health.go
    │       ├── message.go
    │       └── setup.go
    ├── bootstrap                       # This is all Dependency between layers created
    │   └── app.go
    ├── controller                      # Controller layer
    ├── domain                          # Domain Layer
    │   ├── dto
    │   ├── entity
    │   └── interfaces                  # This include interfaces for usecase/controller/repo layers
    ├── repository                      # Repo layer implementation
    ├── usecase                         # Usecase layer implementation
    ├── shared                          # Shared function (logging, etc)
    └── utils                           # Util functions
```

# Installation

There are two ways to set up this project: using Devbox (recommended) or manual installation.

## Option 1: Using Devbox (Recommended)

### Prerequisites
- Devbox: [link](https://www.jetify.com/devbox)

Devbox will automatically install:
- earthly@latest
- go@1.24.3
- trufflehog@latest
- air@latest
- golangci-lint@latest
- psqlodbc@latest

### Steps
1. Start Devbox shell:
```bash
devbox shell
```

2. Start the services:
```bash
devbox run up
```

To shut down:
```bash
devbox run down
```

## Option 2: Manual Installation

### Prerequisites
- Go 1.24.3 or later
- Docker and Docker Compose
- PostgreSQL client (psql)
- Redis CLI (optional, for debugging)

### Steps

1. Install Go dependencies:
```bash
go mod download
```

2. Install development tools:
```bash
# Install air for hot reloading
go install github.com/cosmtrek/air@latest

# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install trufflehog for security scanning
go install github.com/trufflesecurity/trufflehog/v3@latest
```

3. Set up the environment:
```bash
# Copy example environment file
cp .env.example .env

# Edit .env file with your settings
nano .env
```

4. Start the services:
```bash
# For production (with application container)
docker-compose up -d
```

## Configuration

You can configure the application by modifying these environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| MESSAGE_CRON_DURATION | Cron time duration in seconds | 120 |
| MESSAGE_BATCH_NUMBER | Messages handled per batch | 2 |
| WORKER_COUNT | Number of concurrent workers | 2 |
| WORKER_CHAN_BUFFER | Channel buffer size | 100 |
| POSTGRES_HOST | PostgreSQL host | localhost |
| POSTGRES_PORT | PostgreSQL port | 5432 |
| REDIS_HOST | Redis host | localhost |
| REDIS_PORT | Redis port | 6379 |
| WEBHOOK_URL | Webhook URL for sending messages | |
| WEBHOOK_API_KEY | API key for webhook authentication | |

# API Documentation

The application provides REST API endpoints for controlling and monitoring the message processing system. The complete API documentation is available in OpenAPI/Swagger format:

- Swagger JSON: [docs/swagger.json](docs/swagger.json)
- API Documentation: [http://localhost:8080/docs](http://localhost:8080/docs) (when running locally)

Main endpoints include:
- `GET /health` - Health check endpoint
- `POST /service/start` - Start message processing
- `POST /service/stop` - Stop message processing
- `GET /service/status` - Get status of the service
- `GET /messages/sent` - List sent messages

For detailed API documentation including request/response schemas, authentication requirements, and example usage, please refer to the Swagger documentation.

# Development Guide
1. Run docker-compose.dev file
<br>
This file only contains system containers (redis, postgres)

2. Copy .env.example -> .env
<br>
Then config the env webhook

3. Go into Devbox shell (Optional)
<br>
Run ```devbox shell```

4. Enable hot reloading
<br>
Run ```air```

5. Some util commands:
- ```devbox run sql```: This will put you into the postgres container directly so you can do query there.
- ```devbox run redis```: This will put you into the redis container so you can check inside redis with ease.
- ```devbox run lint```: This will execute and check lint for the repo
- ```devbox run scan```: This will scan the repo using Trufflehog for all possible secret leak.
- ```devbox run up```: Deploy the docker-compose
- ```devbox run down```: Tear down the docker-compose resources
- ```devbox run api-gen```: This will generate swagger file from the src code and place it in ***./docs/swagger.json***
- ```devbox run test```: This will run all test cases
 
# Result

### Checklist:
- [x] Retrieve message content, recipient phone number, and sending status for each
record from the Database. Character limit is required for message content.
- [x] Upon project deployment, automatic message sending should -start, processing all
unsent records in the database.
- [x] Messages that have been sent once should not be resent. Newly added records should
be sent in the subsequent automatic process.
- [x] (Bonus Item) - After sending a message, cache to Redis the messageId value received
from the response along with the sending time.

Additional:
- [x] The project should feature 2 separate API endpoints:
    + Start/Stop automatic message sending
    + Retrieve a list of sent messages

Myself-made up:
- [x] Nice Logging
- [x] Can custom duration/batch/maxWorkers
- [x] Follow with clean architecture
- [x] No lint error

>IMPORTANT: Since the webhook.site url provided for me is freetier so I can't set up authentication using API Key header.
However, my code still supports it.
