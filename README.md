# User Service
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com) 
[![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com)

[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ticket-concert_user-service&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ticket-concert_user-service)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ticket-concert_user-service&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ticket-concert_user-service)


## Name
user-service built with :heart:

## Description
User Service is service that used to handle process registration user, login user, update user, view profile, and verify otp registration 

## Installation
1. Ensure, already install golang 1.20 or up
2. Create file .env
```bash
    cp .env.sample .env
```
3. Fill out the env configuration
```bash
#General
SERVICE_NAME=service_user
SERVICE_VERSION=1.0.0
SERVICE_PORT=9001
SERVICE_ENV=development
USERNAME_BASIC_AUTH=username
PASSWORD_BASIC_AUTH=password
SHUTDOWN_DELAY=
SECRET_HASH_PASS=
ID_HASH=

#Mongodb
MONGO_MASTER_DATABASE_URL=mongodb://admin:password@localhost:27020/admin
MONGO_SLAVE_DATABASE_URL=mongodb://admin:password@localhost:27020/admin

#Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

#APM
APM_URL=

#Kafka
KAFKA_URL=localhost:29092

#JWT
JWT_PRIVATE_KEY='your jwt'
JWT_PUBLIC_KEY='your jwt'

JWT_REFRESH_PRIVATE_KEY='your jwt'
JWT_REFRESH_PUBLIC_KEY='your jwt'

APPS_LIMITER=
```
4. Install dependencies:
```bash
make install
```
5. Run in development:
```bash
make run
```

## Test
1. Run unit test
```bash
make unit-test
```
2. Show local coverage (in html)
```bash
make coverage
```

## Authors 
* **Alif Septian Nurdianto** - [Github](https://github.com/alifsn)

## Development Tools
[Fiber](https://gofiber.io/) Rest Framework
[Zap](https://github.com/uber-go/zap) Log Management
[Kafka](https://pkg.go.dev/gopkg.in/confluentinc/confluent-kafka-go.v1@v1.8.2) Event Management
[Mockery](https://github.com/vektra/mockery) Mock Generator
[Go mod](https://go.dev/ref/mod) Depedency Management
[Docker] Container Management