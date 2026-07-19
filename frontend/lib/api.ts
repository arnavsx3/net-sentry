import axios from "axios";

export const api = axios.create({
  baseURL:
    process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api/v1",
  timeout: 10000,
});

export type TargetCurrentAlert = {
  type: string;
  severity: string;
  message: string;
  triggered_at: string;
};

export type TargetCurrentItem = {
  target_host: string;
  status: "healthy" | "degraded" | "down" | "unknown";
  latency_ms: number;
  packet_loss: number;
  observed_at: string | null;
  active_alert_count: number;
  active_alerts: TargetCurrentAlert[];
};

export type CurrentTargetsResponse = {
  count: number;
  targets: TargetCurrentItem[];
};

export type CurrentAlertItem = {
  id: number;
  target_host: string;
  type: string;
  severity: string;
  message: string;
  triggered_at: string;
  resolved_at?: string | null;
  observed_at: string;
  latency_ms: number;
  packet_loss: number;
  probe_status: string;
  probe_result_id: number;
};

export type CurrentAlertsResponse = {
  count: number;
  alerts: CurrentAlertItem[];
};

export type TargetHistoryItem = {
  observed_at: string;
  latency_ms: number;
  packet_loss: number;
  status: "healthy" | "degraded" | "down";
};

export type TargetHistoryResponse = {
  count: number;
  target_host: string;
  results: TargetHistoryItem[];
};

export type TracerouteHop = {
  hop: number;
  address: string;
  rtt_ms: number;
};

export type LatestTracerouteResponse = {
  target_host: string;
  observed_at: string | null;
  probe_status: string;
  latency_ms: number;
  packet_loss: number;
  hops: TracerouteHop[];
};

export type TelemetryEvent = {
  type: "telemetry.received";
  timestamp: string;
  payload: {
    agent_id: string;
    target_host: string;
    observed_at: string;
    status: "healthy" | "degraded" | "down";
    latency_ms: number;
    packet_loss: number;
    trace: TracerouteHop[];
  };
};

export async function getCurrentTargets() {
  const { data } = await api.get<CurrentTargetsResponse>("/targets/current", {
    params: { limit: 50 },
  });
  return data;
}

export async function getCurrentAlerts() {
  const { data } = await api.get<CurrentAlertsResponse>("/alerts/current", {
    params: { limit: 20 },
  });
  return data;
}

export async function getTargetHistory(host: string) {
  const { data } = await api.get<TargetHistoryResponse>(
    `/targets/${encodeURIComponent(host)}/history`,
    { params: { limit: 20 } },
  );
  return data;
}

export async function getLatestTraceroute(host: string) {
  const { data } = await api.get<LatestTracerouteResponse>(
    `/targets/${encodeURIComponent(host)}/traceroute/latest`,
  );
  return data;
}
