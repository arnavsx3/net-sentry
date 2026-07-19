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
