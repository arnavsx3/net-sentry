"use client";

import Link from "next/link";
import { useEffect, useState, useTransition } from "react";

import LiveFeed from "@/app/components/live-feed";
import {
  CurrentAlertItem,
  CurrentAlertsResponse,
  CurrentTargetsResponse,
  TargetCurrentItem,
  TelemetryEvent,
  api,
} from "@/lib/api";

const WS_URL = process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8080/ws";

type Props = {
  initialTargets: CurrentTargetsResponse;
  initialAlerts: CurrentAlertsResponse;
};

function formatDate(value: string | null) {
  if (!value) return "No data yet";
  return new Date(value).toLocaleString();
}

function statusClasses(status: string) {
  switch (status) {
    case "healthy":
      return "bg-emerald-100 text-emerald-700";
    case "degraded":
      return "bg-amber-100 text-amber-700";
    case "down":
      return "bg-rose-100 text-rose-700";
    default:
      return "bg-slate-100 text-slate-600";
  }
}

function severityClasses(severity: string) {
  return severity === "critical"
    ? "bg-rose-100 text-rose-700"
    : "bg-amber-100 text-amber-700";
}

export default function DashboardClient({
  initialTargets,
  initialAlerts,
}: Props) {
  const [targets, setTargets] = useState<TargetCurrentItem[]>(
    initialTargets.targets,
  );
  const [alerts, setAlerts] = useState<CurrentAlertItem[]>(
    initialAlerts.alerts,
  );
  const [, startTransition] = useTransition();

  useEffect(() => {
    const ws = new WebSocket(WS_URL);

    ws.onmessage = (event) => {
      try {
        const parsed = JSON.parse(event.data) as
          | TelemetryEvent
          | { type: string };

        if (parsed.type !== "telemetry.received") {
          return;
        }

        const telemetry = parsed as TelemetryEvent;

        setTargets((current) => {
          const next = [...current];
          const index = next.findIndex(
            (item) => item.target_host === telemetry.payload.target_host,
          );

          const updatedTarget: TargetCurrentItem = {
            target_host: telemetry.payload.target_host,
            status: telemetry.payload.status,
            latency_ms: telemetry.payload.latency_ms,
            packet_loss: telemetry.payload.packet_loss,
            observed_at: telemetry.payload.observed_at,
            active_alert_count: index >= 0 ? next[index].active_alert_count : 0,
            active_alerts: index >= 0 ? next[index].active_alerts : [],
          };

          if (index >= 0) {
            next[index] = updatedTarget;
          } else {
            next.push(updatedTarget);
          }

          next.sort((a, b) => a.target_host.localeCompare(b.target_host));
          return next;
        });

        startTransition(async () => {
          try {
            const { data } = await api.get<CurrentAlertsResponse>(
              "/alerts/current",
              {
                params: { limit: 20 },
              },
            );

            setAlerts(data.alerts);

            setTargets((current) =>
              current.map((target) => {
                const targetAlerts = data.alerts.filter(
                  (alert) => alert.target_host === target.target_host,
                );

                return {
                  ...target,
                  active_alert_count: targetAlerts.length,
                  active_alerts: targetAlerts.map((alert) => ({
                    type: alert.type,
                    severity: alert.severity,
                    message: alert.message,
                    triggered_at: alert.triggered_at,
                  })),
                };
              }),
            );
          } catch {
            // ignore transient refresh failures for now
          }
        });
      } catch {
        // ignore malformed messages
      }
    };

    return () => {
      ws.close();
    };
  }, []);

  return (
    <main className="min-h-screen bg-slate-950 text-slate-100">
      <div className="mx-auto flex max-w-7xl flex-col gap-10 px-6 py-10">
        <section className="space-y-4">
          <p className="text-sm font-medium uppercase tracking-[0.24em] text-cyan-400">
            NetSentry
          </p>
          <div className="space-y-3">
            <h1 className="max-w-4xl text-4xl font-semibold tracking-tight text-white sm:text-5xl">
              Network observability dashboard for live target health and alert
              tracking.
            </h1>
            <p className="max-w-2xl text-base leading-7 text-slate-400">
              Current target status, active threshold alerts, and latest
              telemetry from your Go backend.
            </p>
          </div>
        </section>

        <section className="grid gap-4 sm:grid-cols-2">
          <div className="rounded-3xl border border-slate-800 bg-slate-900 p-6 shadow-2xl shadow-black/20">
            <p className="text-sm text-slate-400">Targets</p>
            <h2 className="mt-3 text-4xl font-semibold text-white">
              {targets.length}
            </h2>
          </div>

          <div className="rounded-3xl border border-slate-800 bg-slate-900 p-6 shadow-2xl shadow-black/20">
            <p className="text-sm text-slate-400">Active Alerts</p>
            <h2 className="mt-3 text-4xl font-semibold text-white">
              {alerts.length}
            </h2>
          </div>
        </section>

        <section className="space-y-5">
          <h2 className="text-2xl font-semibold text-white">Current Targets</h2>

          <div className="grid gap-5 lg:grid-cols-2 xl:grid-cols-3">
            {targets.map((target) => (
              <Link
                key={target.target_host}
                href={`/targets/${encodeURIComponent(target.target_host)}`}
                className="rounded-3xl border border-slate-800 bg-slate-900 p-6 shadow-2xl shadow-black/20 transition hover:border-cyan-500/50 hover:bg-slate-900/90">
                <div className="mb-5 flex items-start justify-between gap-4">
                  <div>
                    <p className="text-xs uppercase tracking-[0.2em] text-slate-500">
                      Target
                    </p>
                    <h3 className="mt-2 text-2xl font-semibold text-white">
                      {target.target_host}
                    </h3>
                  </div>

                  <span
                    className={`rounded-full px-3 py-1 text-xs font-semibold uppercase tracking-[0.14em] ${statusClasses(
                      target.status,
                    )}`}>
                    {target.status}
                  </span>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div className="rounded-2xl bg-slate-800/70 p-4">
                    <p className="text-xs uppercase tracking-[0.14em] text-slate-400">
                      Latency
                    </p>
                    <p className="mt-2 text-2xl font-semibold text-white">
                      {target.latency_ms} ms
                    </p>
                  </div>

                  <div className="rounded-2xl bg-slate-800/70 p-4">
                    <p className="text-xs uppercase tracking-[0.14em] text-slate-400">
                      Packet Loss
                    </p>
                    <p className="mt-2 text-2xl font-semibold text-white">
                      {target.packet_loss}%
                    </p>
                  </div>
                </div>

                <div className="mt-5 space-y-2 text-sm text-slate-400">
                  <p>Last observed: {formatDate(target.observed_at)}</p>
                  <p>Active alerts: {target.active_alert_count}</p>
                </div>

                <div className="mt-5 space-y-3">
                  {target.active_alerts.length > 0 ? (
                    target.active_alerts.map((alert, index) => (
                      <div
                        key={`${target.target_host}-${alert.type}-${index}`}
                        className="rounded-2xl border border-slate-800 bg-slate-800/50 p-4">
                        <p className="text-sm font-semibold text-white">
                          {alert.type} - {alert.severity}
                        </p>
                        <p className="mt-2 text-sm leading-6 text-slate-400">
                          {alert.message}
                        </p>
                      </div>
                    ))
                  ) : (
                    <div className="rounded-2xl border border-dashed border-slate-700 p-4 text-sm text-slate-500">
                      No active alerts.
                    </div>
                  )}
                </div>
              </Link>
            ))}
          </div>
        </section>

        <LiveFeed />

        <section className="space-y-5">
          <h2 className="text-2xl font-semibold text-white">Current Alerts</h2>

          <div className="overflow-hidden rounded-3xl border border-slate-800 bg-slate-900 shadow-2xl shadow-black/20">
            {alerts.length === 0 ? (
              <div className="p-6 text-sm text-slate-400">
                No active alerts.
              </div>
            ) : (
              <div className="divide-y divide-slate-800">
                {alerts.map((alert) => (
                  <div key={alert.id} className="p-6">
                    <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                      <div>
                        <h3 className="text-lg font-semibold text-white">
                          {alert.target_host} - {alert.type}
                        </h3>
                        <p className="mt-2 text-sm leading-6 text-slate-400">
                          {alert.message}
                        </p>
                      </div>

                      <span
                        className={`w-fit rounded-full px-3 py-1 text-xs font-semibold uppercase tracking-[0.14em] ${severityClasses(
                          alert.severity,
                        )}`}>
                        {alert.severity}
                      </span>
                    </div>

                    <div className="mt-4 grid gap-3 text-sm text-slate-400 sm:grid-cols-3">
                      <p>Observed: {formatDate(alert.observed_at)}</p>
                      <p>Latency: {alert.latency_ms} ms</p>
                      <p>Packet loss: {alert.packet_loss}%</p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </section>
      </div>
    </main>
  );
}