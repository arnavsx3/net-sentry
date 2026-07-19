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
