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