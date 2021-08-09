# NotifyHub

Simple web server for sending user notifications through some communication channels:

 - Email
 - SMS
 - Telegram

## Techical dependencies

 - RabbitMQ message broker
 - PostgreSQL database management system

## Features

 - Multithreaded queue processing
 - PostgreSQL DB for storing recipient data and logging
 - API for creating recipients and sending notifications

## Quick Start

### General usage

1) Start RabbitMQ consumer
```bash
go run cmd/sender/sender.go
```

2) Start API web server
```bash
go run api/server.go --port=5000
```

3) Create recipient
  Send POST http request on /recipients
```json
{
  "email": "example@email.com",
  "phone": "phone number with + sign",
  "telegram": "telegram chat id"
}
```
4) Send notification
  Send POST http request on /submit
```json
{
  "notificationTypes": ["email", "telegram"],
  "recipientList": [1, 2, 3],
  "message": "message text"
}
```