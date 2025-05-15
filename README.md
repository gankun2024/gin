---
title: Gin
description: A Gin server
tags:
  - gin
  - golang
---

# Go Gin API Service

A RESTful API service built with Go and the Gin framework, featuring authentication and Stripe payment integration.

## Directory Structure

```
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── routes/
│   ├── config/
│   ├── db/
│   │   ├── migrations/
│   │   └── models/
│   └── services/
│       ├── auth/
│       └── payment/
├── pkg/
│   ├── logger/
│   └── utils/
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Features

- User authentication (register, login, password reset)
- JWT-based authentication middleware
- Stripe payment integration
- Checkout session creation
- Webhook handling for payment events
- Environment-based configuration
- Structured logging

## Getting Started

### Prerequisites

- Go 1.22 or higher
- PostgreSQL (optional, for database integration)
- Stripe account for payments

### Installation

1. Clone the repository:
```bash
git clone https://github.com/gankun2024/gin-demo-project.git
cd yourproject
```

2. Copy the example environment file and configure it:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Run the application:
```bash
go run cmd/api/main.go
```

### Environment Variables

Configure the following environment variables in your `.env` file:

- `PORT`: Server port (default: 8080)
- `GIN_MODE`: Gin mode (debug, release, test)
- `LOG_LEVEL`: Log level (debug, info, warn, error)
- `DB_*`: Database connection settings
- `JWT_SECRET`: Secret key for JWT tokens
- `STRIPE_*`: Stripe API keys and webhook secret

## API Endpoints

### Authentication

- `POST /api/v1/auth/register`: Register a new user
- `POST /api/v1/auth/login`: Authenticate a user
- `POST /api/v1/auth/forgot-password`: Request password reset
- `POST /api/v1/auth/reset-password`: Reset password with token

### User

- `GET /api/v1/user/profile`: Get user profile
- `PUT /api/v1/user/profile`: Update user profile

### Payments

- `POST /api/v1/payments/create-checkout`: Create Stripe checkout session
- `GET /api/v1/payments/success`: Handle successful payment
- `GET /api/v1/payments/cancel`: Handle canceled payment

### Webhooks

- `POST /webhooks/stripe`: Handle Stripe webhook events

## Development

### Adding a New API Endpoint

1. Create a handler function in the appropriate file under `internal/api/handlers/`
2. Register the endpoint in `internal/api/routes/routes.go`

### Database Migrations

Database migrations would be added in the `internal/db/migrations/` directory.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
