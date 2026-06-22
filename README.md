# NetSentry

NetSentry is a distributed network observability platform built to monitor network health, detect anomalies, visualize routing paths, and provide real-time operational insight into latency, packet loss, and route changes.

## Goals

NetSentry is designed as a production-grade resume project that demonstrates:

- distributed systems design
- Go backend development
- network telemetry collection
- time-series data handling
- alerting and anomaly detection
- real-time dashboards
- path and topology visualization

## Core Features

- Distributed Go agents that run active probes every 30 seconds
- Latency and packet loss measurement across configured targets
- Traceroute collection and route change detection
- Historical telemetry storage in PostgreSQL
- Real-time alerting for unhealthy links
- Basic anomaly detection using statistical baselines
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

### Observability
- Prometheus
- Grafana

### Deployment
- Docker
- Docker Compose

## Repository Structure

```text
agent/      Go-based network probe agent
backend/    Gin API, ingestion, analysis, alerting
frontend/   Next.js dashboard
deploy/     Docker Compose and deployment configs
docs/       Architecture and scope documents