import Link from "next/link";
import { redirect } from "next/navigation";

import { AppShell } from "@/components/app-shell";
import { Card } from "@/components/ui/card";
import { requireSession } from "@/lib/auth";
import { listPlans } from "@/lib/data";

export const dynamic = "force-dynamic";

export default async function PlansPage() {
  try {
    const session = await requireSession();
    const plans = await listPlans(session.user.id);

    return (
      <AppShell>
        <Card className="space-y-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs uppercase tracking-[0.2em] text-black/45">
                Stream plans
              </p>
              <h2 className="text-2xl font-semibold">Planning queue</h2>
            </div>
            <Link className="text-sm text-ember underline" href="/plans/new">
              Create plan
            </Link>
          </div>
          <div className="space-y-4">
            {plans.length === 0 ? (
              <p className="text-sm text-black/55">No plans created yet.</p>
            ) : (
              plans.map((plan) => (
                <Link
                  key={plan.id}
                  href={`/plans/${plan.id}`}
                  className="block rounded-2xl border border-black/10 bg-sand/70 px-4 py-3"
                >
                  <div className="flex items-center justify-between gap-4">
                    <div>
                      <p className="font-medium">{plan.name}</p>
                      <p className="text-sm text-black/55">
                        {String(plan.plan_date).slice(0, 10)} • {plan.status}
                      </p>
                    </div>
                    <span className="text-xs uppercase tracking-[0.18em] text-black/40">
                      {plan.auto_title_enabled ? "Auto titles" : "Manual"}
                    </span>
                  </div>
                </Link>
              ))
            )}
          </div>
        </Card>
      </AppShell>
    );
  } catch {
    redirect("/login");
  }
}
