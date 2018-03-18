# echo
A simple web service that provides a few basic endpoints.

[![CircleCI](https://circleci.com/gh/blupig/echo.svg?style=svg)](https://circleci.com/gh/blupig/echo)
[![Coverage Status](https://coveralls.io/repos/github/blupig/echo/badge.svg?branch=ci)](https://coveralls.io/github/blupig/echo?branch=ci)

## Endpoints
API is avaliable at https://api.blupig.net/echo

Avaliable endpoints:
- `GET /cache`: cacheable content with increased latency (help debugging caching layer)
- `GET /cpu`: performs CPU-intensive operation on server-side (requires API token)
- `GET /exit`: causes server process to exit (requires API token)
- `GET /headers`: returns request headers in JSON
- `GET /health`: returns application health info
- `GET /ip`: returns client IP (uses `X-Forwarded-For` if present, otherwise returns client IP)

## Deploy
Pre-built binaries are available as Docker images at `blupig/echo`.

The server can be configured with environment variables:
- `PORT`: the port server listens on (default: `8000`)
- `API_TOKEN`: API token for `/cpu` and `/exit` endpoints, if not set or set to empty string, all endpoints require API token are disabled.
