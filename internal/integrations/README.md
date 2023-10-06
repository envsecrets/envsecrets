# Integrations Service

This service handles connections and communication with all the third-party platforms like Github and Vercel with whom we want to sync our secrets with.

Note: Always abide by clean code architecture and implement proper abstractions between different layers of business logic.

## Style Guide

Mandatory files:

- `routes.go` contains API routes.
- `handler.go` contains handlers for every route.
- `default.go` sets and gets the default initialized service.
- `init.go` initializes the default service at startup.

## Contribution Guide

To add a new integration:

1. Register it's unique `IntegrationType` constant in `commons.go`. For example: `const Github IntegrationType = "github"`
1. 
