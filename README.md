# Overview
This project is interview test from Insider for Senior Golang Developer role

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
├── .golangci.yaml
├── .air.toml
├── .env.example
├── Dockerfile
├── Dockerfile
├── LICENSE
├── README.md
├── build
│   └── init.sql
├── cmd
│   └── server
│       └── main.go
├── devbox.json
├── devbox.lock
├── docker-compose.yml
├── docs
├── e2e
│   └── main_test.go
├── internal
    ├── api
    │   ├── middleware
    │   │   └── logging.go
    │   └── route
    │       ├── health.go
    │       ├── message.go
    │       └── setup.go
    ├── bootstrap
    │   └── app.go
    ├── controller
    ├── domain
    │   ├── dto
    │   ├── entity
    │   └── interfaces
    ├── repository
    ├── shared
    ├── usecase
    └── utils
```

# How to setup
This section guide how to setup the repository. Personally, I like to use Devbox - which is similar with package.json for javascript. But Devbox supports to have all required libraries/binaries for your dev env as well as scripting utilities. You can use or not, please looks into ***./devbox.json*** file for all util script.

## Prerequisites
- Devbox: [link](https://www.jetify.com/devbox)

Devbox will installs:
- earthly@latest
- go@1.24.3
- trufflehog@latest
- air@latest
- golangci-lint@latest
- psqlodbc@latest

## Guidelines
The project and all of it's system is setup using docker-compose.

1. If you are using devbox then it's quite simple
```
devbox shell
devbox up
```

>***devbox shell*** will create a shell where it has all required libraries with the exact version declared in devbox.json

>***devbox up*** will setup docker-compose

2. If you want to do it your self:
<br>
Look into ***devbox.json*** for the required command
