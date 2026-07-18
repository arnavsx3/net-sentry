const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api/v1";

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