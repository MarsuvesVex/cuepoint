import Link from "next/link";
import { redirect } from "next/navigation";

import { AppShell } from "@/components/app-shell";
import { Card } from "@/components/ui/card";
import { requireSession } from "@/lib/auth";
import { getDashboardData } from "@/lib/data";

export const dynamic = "force-dynamic";

export default async function DashboardPage() {
  try {
    const session = await requireSession();
    const data = await getDashboardData(session.user.id);

    return (
      <AppShell>
        <div className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
          <Card className="space-y-6">
            <div className="flex items-start justify-between gap-4">
              <div>
                <p className="text-xs uppercase tracking-[0.2em] text-black/45 dark:text-white/45">
                  Overview
                </p>
                <h2 className="text-2xl font-semibold">
                  {data.profile?.twitch_display_name || session.user.name || "Operator"}
                </h2>
                <p className="mt-2 text-sm text-black/60 dark:text-white/60">
                  Use the dashboard to watch current runtime defaults and jump into plans or Twitch settings.
                </p>
              </div>
              <Link className="text-sm text-ember underline" href="/plans/new">
                New plan
              </Link>
            </div>
            <div className="grid gap-4 md:grid-cols-3">
              <Metric label="Linked Twitch" value={data.profile?.twitch_login || "Not linked"} />
              <Metric label="Recent plans" value={String(data.plans.length)} />
              <Metric
                label="Live title"
                value={data.liveSession?.current_title || "Offline"}
              />
            </div>
            <div className="space-y-3">
              <h3 className="text-sm font-medium text-black/70 dark:text-white/70">Recent plans</h3>
              {data.plans.length === 0 ? (
                <p className="text-sm text-black/55 dark:text-white/55">No plans yet.</p>
              ) : (
                <div className="space-y-3">
                  {data.plans.map((plan) => (
                    <Link
                      key={plan.id}
                      href={`/plans/${plan.id}`}
                      className="block rounded-2xl border border-black/10 bg-sand/70 px-4 py-3 dark:border-white/10 dark:bg-white/5"
                    >
                      <div className="flex items-center justify-between gap-4">
                        <div>
                          <p className="font-medium">{plan.name}</p>
                          <p className="text-sm text-black/55 dark:text-white/55">
                            {String(plan.plan_date).slice(0, 10)} • {plan.status}
                          </p>
                        </div>
                        <span className="text-xs uppercase tracking-[0.18em] text-black/40 dark:text-white/40">
                          {plan.auto_title_enabled ? "Auto titles on" : "Manual"}
                        </span>
                      </div>
                    </Link>
                  ))}
                </div>
              )}
            </div>
          </Card>
          <div className="space-y-6">
            <Card className="space-y-4">
              <h3 className="text-lg font-semibold">Live session</h3>
              <p className="text-sm text-black/65 dark:text-white/65">
                {data.liveSession
                  ? `Current title: ${data.liveSession.current_title}`
                  : "No active stream session yet. The bot will create one after sync when you go live."}
              </p>
              {data.devSettings?.force_live_state ? (
                <p className="rounded-2xl bg-ember/10 px-4 py-3 text-sm text-ember dark:bg-ember/20">
                  Dev live override is enabled. Runtime stream checks will report a forced live state.
                </p>
              ) : null}
            </Card>
            <Card className="space-y-4">
              <h3 className="text-lg font-semibold">Bot defaults</h3>
              <p className="text-sm text-black/65 dark:text-white/65">
                {data.settings?.default_title_format || "{original_title} | {segment_title}"}
              </p>
              <Link className="text-sm text-ember underline" href="/settings/bot">
                Open bot settings
              </Link>
            </Card>
          </div>
        </div>
      </AppShell>
    );
  } catch {
    redirect("/login");
  }
}

function Metric({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-2xl bg-sand/70 px-4 py-3 dark:bg-white/5">
      <p className="text-xs uppercase tracking-[0.18em] text-black/40 dark:text-white/40">{label}</p>
      <p className="mt-2 text-sm font-medium">{value}</p>
    </div>
  );
}
