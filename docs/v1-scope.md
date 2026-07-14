# NetSentry V1 Scope

### Agent
- Written in Go
- Runs every 30 seconds
- Executes:
  - ping
  - packet loss measurement
  - traceroute
- Sends telemetry to backend over HTTP
- Supports retries on transient failures
- Supports config-driven target lists

### Backend
- Written in Go with Gin
- Exposes ingestion APIs for telemetry submission
- Validates and stores agent data
- Persists historical measurements in PostgreSQL
- Detects:
  - high latency
  - packet loss threshold breaches
- Exposes APIs for dashboard queries
- Pushes live events over WebSockets

### Frontend
- Built with Next.js and TypeScript
- Displays:
  - monitored targets
  - current network health
  - latency trends
  - packet loss trends
  - active alerts
  - traceroute snapshots

### Deployment
- Dockerized services
- Local orchestration with Docker Compose
- Basic service health checks

## Out of Scope for V1

These are intentionally excluded unless extra time remains:

- full enterprise topology auto-discovery
- SNMP-based discovery
- deep bandwidth analytics
- route change detection automation
- z-score anomaly detection
- RBAC and multi-user auth
- Kubernetes deployment
- complex failure simulation engine
- Prometheus and Grafana integration

## Success Criteria

NetSentry V1 is complete when:

- at least one agent can send telemetry successfully
- backend stores and serves historical data
- alerts are generated for unhealthy links
- frontend shows live and historical network health
- traceroute snapshots can be viewed
- the full system runs locally with Docker Compose

## Non-Functional Goals

- clear code structure
- readable API contracts
- basic logging and error handling
- modular services
- good resume-level architecture
