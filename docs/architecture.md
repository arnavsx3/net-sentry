# NetSentry Architecture

## Overview

NetSentry follows a distributed telemetry collection model.

Go-based agents run active network probes from different locations and send results to a centralized backend. The backend stores telemetry, analyzes health patterns, generates alerts, and serves real-time and historical data to a web dashboard.

## High-Level Components

### 1. Agent

The agent runs on remote machines or local nodes.

Responsibilities:

- execute ping probes
- measure latency
- estimate packet loss
- run traceroute
- package telemetry results
- send results to backend on a fixed schedule

### 2. Backend API

The backend is the central ingestion and query service.

Responsibilities:

- receive telemetry from agents
- validate request payloads
- persist measurements to PostgreSQL
- expose APIs for dashboard queries
- publish live updates through WebSockets

### 3. Analysis and Alerting

This layer evaluates incoming telemetry.

Responsibilities:

- detect latency threshold breaches
- detect packet loss threshold breaches
- detect route changes between snapshots
- compute rolling baselines
- flag anomalies using z-score
- generate alert records

### 4. Database

PostgreSQL stores operational and historical data.

Expected entities:

- agents
- targets
- probe_results
- traceroute_hops
- alerts
- incidents

### 5. Frontend Dashboard

The dashboard is the operator-facing UI.

Responsibilities:

- show current target health
- render latency and packet loss charts
- display alert feed
- visualize traceroute paths
- show incident and recovery history

### 6. Platform Observability

NetSentry should observe itself as well.

Responsibilities:

- expose backend metrics for Prometheus
- visualize internal service metrics in Grafana
- surface service health status

## Request Flow

1. Agent loads configured targets
2. Agent runs ping and traceroute checks
3. Agent sends telemetry payload to backend
4. Backend stores raw results in PostgreSQL
5. Analysis logic evaluates thresholds and anomalies
6. Alerts and route changes are recorded
7. Frontend fetches history and receives live updates via WebSockets
