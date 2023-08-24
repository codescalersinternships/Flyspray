# Flyspray

Flyspray is a web-based bug tracking system.

## Table of Contents

- [Introduction](#introduction)
- [Configuration](#Configuration)
- [Installation](#installation)
- [Testing](#testing)

## Introduction

FlySpray is a web-based application designed to manage bug tracking for software projects. The system provides tools for users to report, monitor, and resolve bugs throughout the software development lifecycle. It offers features like user authentication, project management, bug creation and tracking, as well as user comments on bugs.

## Configuration

Before building or running backend image create `config.json` in `backend` dir.

`config.json` example:
```json
{
  "server": {
    "host": "localhost",
    "port": 8080
  },
  "mail_sender": {
    "email": "<email>",
    "sendgrid_key": "<sendgrid-key>",
    "timeout": 30
  },
  "db": {
    "file": "./database.db"
  },
  "jwt": {
    "secret": "<secret>",
    "timeout": 5
  }
}
```

## Installation

Make sure you have added `config.json` (see [Configuration](#Configuration)) and have installed the following tools before proceeding:

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/install/)

To quickly set up and run the FlySpray, follow these steps:

1. Clone this repository to your local machine:

```bash
git clone https://github.com/codescalersinternships/Flyspray
cd Flyspray
```

2. Start the FlySpray:

```bash
make run
```
## Testing

First install all dependencies and you can test frontend and backend

1. install backend dependencies
```bash
cd backend
go mod download
```

2. install frontend dependencies
```bash
cd frontend
npm install
```

3. test backend and frontend
```bash
make test
```