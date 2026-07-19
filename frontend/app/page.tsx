import DashboardClient from "@/app/components/dashboard-client";
import { getCurrentAlerts, getCurrentTargets } from "@/lib/api";

export default async function HomePage() {
  const [targetsData, alertsData] = await Promise.all([
    getCurrentTargets(),
    getCurrentAlerts(),
  ]);

  return (
    <DashboardClient initialTargets={targetsData} initialAlerts={alertsData} />
  );
}
