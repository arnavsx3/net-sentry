import Link from "next/link";
import { notFound } from "next/navigation";
import { getLatestTraceroute, getTargetHistory } from "@/lib/api";
import HistoryChart from "@/app/components/history-chart";

type Props = {
  params: Promise<{
    host: string;
  }>;
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

export default async function TargetDetailsPage({ params }: Props) {
  const { host } = await params;
  const decodedHost = decodeURIComponent(host);

  const [historyData, tracerouteData] = await Promise.all([
    getTargetHistory(decodedHost),
    getLatestTraceroute(decodedHost),
  ]);

  if (!historyData.target_host && !tracerouteData.target_host) {
    notFound();
  }

  return (
    <main className="min-h-screen bg-slate-950 text-slate-100">
      <div className="mx-auto flex max-w-6xl flex-col gap-8 px-6 py-10">
        <div className="flex items-center justify-between">
          <div>
            <Link
              href="/"
              className="text-sm text-cyan-400 hover:text-cyan-300">
              Back to dashboard
            </Link>
            <h1 className="mt-3 text-4xl font-semibold text-white">
              {decodedHost}
            </h1>
            <p className="mt-2 text-slate-400">
              Detailed telemetry history and latest traceroute snapshot.
            </p>
          </div>

          <span
            className={`rounded-full px-3 py-1 text-xs font-semibold uppercase tracking-[0.14em] ${statusClasses(
              tracerouteData.probe_status || "unknown",
            )}`}>
            {tracerouteData.probe_status || "unknown"}
          </span>
        </div>

        <section className="grid gap-4 md:grid-cols-3">
          <div className="rounded-3xl border border-slate-800 bg-slate-900 p-6">
            <p className="text-sm text-slate-400">Latest Latency</p>
            <h2 className="mt-2 text-3xl font-semibold text-white">
              {tracerouteData.latency_ms} ms
            </h2>
          </div>

          <div className="rounded-3xl border border-slate-800 bg-slate-900 p-6">
            <p className="text-sm text-slate-400">Latest Packet Loss</p>
            <h2 className="mt-2 text-3xl font-semibold text-white">
              {tracerouteData.packet_loss}%
            </h2>
          </div>

          <div className="rounded-3xl border border-slate-800 bg-slate-900 p-6">
            <p className="text-sm text-slate-400">Last Observed</p>
            <h2 className="mt-2 text-lg font-semibold text-white">
              {formatDate(tracerouteData.observed_at)}
            </h2>
          </div>
        </section>

        <section className="space-y-5">
          <h2 className="text-2xl font-semibold text-white">
            Latency and Packet Loss Trend
          </h2>
          {historyData.results.length === 0 ? (
            <div className="rounded-3xl border border-slate-800 bg-slate-900 p-6 text-slate-400">
              No history available yet.
            </div>
          ) : (
            <HistoryChart data={historyData.results} />
          )}
        </section>

        <section className="space-y-5">
          <h2 className="text-2xl font-semibold text-white">Recent History</h2>

          <div className="overflow-hidden rounded-3xl border border-slate-800 bg-slate-900">
            {historyData.results.length === 0 ? (
              <div className="p-6 text-slate-400">
                No history available yet.
              </div>
            ) : (
              <div className="divide-y divide-slate-800">
                {historyData.results.map((item, index) => (
                  <div
                    key={`${item.observed_at}-${index}`}
                    className="grid gap-4 p-6 md:grid-cols-4 md:items-center">
                    <div>
                      <p className="text-xs uppercase tracking-[0.14em] text-slate-500">
                        Observed
                      </p>
                      <p className="mt-2 text-sm text-slate-300">
                        {formatDate(item.observed_at)}
                      </p>
                    </div>

                    <div>
                      <p className="text-xs uppercase tracking-[0.14em] text-slate-500">
                        Latency
                      </p>
                      <p className="mt-2 text-lg font-semibold text-white">
                        {item.latency_ms} ms
                      </p>
                    </div>

                    <div>
                      <p className="text-xs uppercase tracking-[0.14em] text-slate-500">
                        Packet Loss
                      </p>
                      <p className="mt-2 text-lg font-semibold text-white">
                        {item.packet_loss}%
                      </p>
                    </div>

                    <div>
                      <p className="text-xs uppercase tracking-[0.14em] text-slate-500">
                        Status
                      </p>
                      <span
                        className={`mt-2 inline-flex rounded-full px-3 py-1 text-xs font-semibold uppercase tracking-[0.14em] ${statusClasses(
                          item.status,
                        )}`}>
                        {item.status}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </section>

        <section className="space-y-5">
          <h2 className="text-2xl font-semibold text-white">
            Latest Traceroute
          </h2>

          <div className="overflow-hidden rounded-3xl border border-slate-800 bg-slate-900">
            {tracerouteData.hops.length === 0 ? (
              <div className="p-6 text-slate-400">
                No traceroute hops in the latest sample.
              </div>
            ) : (
              <div className="divide-y divide-slate-800">
                {tracerouteData.hops.map((hop) => (
                  <div
                    key={`${hop.hop}-${hop.address}`}
                    className="grid gap-4 p-6 md:grid-cols-3 md:items-center">
                    <div>
                      <p className="text-xs uppercase tracking-[0.14em] text-slate-500">
                        Hop
                      </p>
                      <p className="mt-2 text-lg font-semibold text-white">
                        {hop.hop}
                      </p>
                    </div>

                    <div>
                      <p className="text-xs uppercase tracking-[0.14em] text-slate-500">
                        Address
                      </p>
                      <p className="mt-2 text-sm text-slate-300">
                        {hop.address}
                      </p>
                    </div>

                    <div>
                      <p className="text-xs uppercase tracking-[0.14em] text-slate-500">
                        RTT
                      </p>
                      <p className="mt-2 text-lg font-semibold text-white">
                        {hop.rtt_ms} ms
                      </p>
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
