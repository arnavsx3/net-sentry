# NetSentry Architecture

## Overview

NetSentry follows a simple telemetry collection model for a strong MVP.

Go-based agents run active network probes on a schedule and send results to a centralized backend. The backend stores telemetry, applies threshold-based health checks, and serves live and historical data to a web dashboard.

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

This layer evaluates incoming telemetry with simple rules.

Responsibilities:

- detect latency threshold breaches
- detect packet loss threshold breaches
- generate alert records

### 4. Database

PostgreSQL stores operational and historical data.

Expected entities:

- agents
- targets
- probe_results
- traceroute_hops
- alerts

### 5. Frontend Dashboard

The dashboard is the operator-facing UI.

Responsibilities:

- show current target health
- render latency and packet loss charts
- display alert feed
- show traceroute snapshots and recent paths

## Request Flow

1. Agent loads configured targets
2. Agent runs ping and traceroute checks
3. Agent sends telemetry payload to backend
4. Backend stores raw results in PostgreSQL
5. Analysis logic evaluates threshold breaches
6. Alerts are recorded
7. Frontend fetches history and receives live updates via WebSockets
