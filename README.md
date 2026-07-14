# NetSentry

NetSentry is a network observability MVP built to monitor network health, collect active probe telemetry, and provide real-time insight into latency, packet loss, and traceroute results.

## Goals

NetSentry is designed as a resume-focused project that demonstrates:

- distributed systems design
- Go backend development
- network telemetry collection
- time-series data handling
- threshold-based alerting
- real-time dashboards
- traceroute inspection

## Core Features

- Go agents that run active probes every 30 seconds
- Latency and packet loss measurement across configured targets
- Traceroute collection for route inspection
- Historical telemetry storage in PostgreSQL
- Real-time alerting for unhealthy links
- Live updates over WebSockets
- Live dashboard built with Next.js and TypeScript
- Local deployment with Docker Compose

## Tech Stack

### Backend
- Go
- Gin

### Database
- PostgreSQL

### Frontend
- Next.js
- TypeScript

### Realtime
- WebSockets

### Deployment
- Docker
- Docker Compose

## Repository Structure

```text
agent/      Go-based network probe agent
backend/    Gin API, ingestion, storage, alerting
frontend/   Next.js dashboard
deploy/     Docker Compose and deployment configs
docs/       Architecture and scope documents
