"use client";

import { useEffect, useState } from "react";

type TraceHop = {
  hop: number;
  address: string;
  rtt_ms: number;
};

type TelemetryEvent = {
  type: string;
  timestamp: string;
  payload: {
    agent_id: string;
    target_host: string;
    observed_at: string;
    status: "healthy" | "degraded" | "down";
    latency_ms: number;
    packet_loss: number;
    trace: TraceHop[];
  };
};

const WS_URL = process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8080/ws";

function formatDate(value: string) {
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

export default function LiveFeed() {
  const [connected, setConnected] = useState(false);
  const [events, setEvents] = useState<TelemetryEvent[]>([]);

  useEffect(() => {
    const ws = new WebSocket(WS_URL);

    ws.onopen = () => {
      setConnected(true);
    };

    ws.onclose = () => {
      setConnected(false);
    };

    ws.onerror = () => {
      setConnected(false);
    };

    ws.onmessage = (event) => {
      try {
        const parsed = JSON.parse(event.data) as
          | TelemetryEvent
          | { type: string };

        if (parsed.type !== "telemetry.received") {
          return;
        }

        setEvents((current) =>
          [parsed as TelemetryEvent, ...current].slice(0, 10),
        );
      } catch {
        // ignore malformed payloads for now
      }
    };

    return () => {
      ws.close();
    };
  }, []);

  return (
    <section className="space-y-5">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-semibold text-white">Live Feed</h2>
        <span
          className={`rounded-full px-3 py-1 text-xs font-semibold uppercase tracking-[0.14em] ${
            connected
              ? "bg-emerald-100 text-emerald-700"
              : "bg-slate-200 text-slate-700"
          }`}>
          {connected ? "connected" : "disconnected"}
        </span>
      </div>

      <div className="overflow-hidden rounded-3xl border border-slate-800 bg-slate-900 shadow-2xl shadow-black/20">
        {events.length === 0 ? (
          <div className="p-6 text-sm text-slate-400">
            Waiting for live telemetry events...
          </div>
        ) : (
          <div className="divide-y divide-slate-800">
            {events.map((item, index) => (
              <div key={`${item.timestamp}-${index}`} className="p-6">
                <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                  <div>
                    <h3 className="text-lg font-semibold text-white">
                      {item.payload.target_host} - {item.payload.status}
                    </h3>
                    <p className="mt-2 text-sm leading-6 text-slate-400">
                      Agent: {item.payload.agent_id}
                    </p>
                  </div>

                  <span
                    className={`w-fit rounded-full px-3 py-1 text-xs font-semibold uppercase tracking-[0.14em] ${statusClasses(
                      item.payload.status,
                    )}`}>
                    {item.payload.status}
                  </span>
                </div>

                <div className="mt-4 grid gap-3 text-sm text-slate-400 sm:grid-cols-4">
                  <p>Observed: {formatDate(item.payload.observed_at)}</p>
                  <p>Latency: {item.payload.latency_ms} ms</p>
                  <p>Packet loss: {item.payload.packet_loss}%</p>
                  <p>Trace hops: {item.payload.trace.length}</p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </section>
  );
}
