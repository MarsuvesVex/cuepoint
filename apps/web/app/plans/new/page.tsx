import { redirect } from "next/navigation";

import { AppShell } from "@/components/app-shell";
import { FormField } from "@/components/form-field";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Select } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { requireSession } from "@/lib/auth";
import { createPlanAction } from "@/lib/actions";

export const dynamic = "force-dynamic";

export default async function NewPlanPage() {
  try {
    await requireSession();

    return (
      <AppShell>
        <Card className="max-w-3xl space-y-6">
          <div>
            <p className="text-xs uppercase tracking-[0.2em] text-black/45 dark:text-white/45">
              New plan
            </p>
            <h2 className="text-2xl font-semibold">Create a stream plan</h2>
            <p className="mt-2 text-sm text-black/60 dark:text-white/60">
              Set the date, scheduling mode, and default automation behavior for a single stream run.
            </p>
          </div>
          <form action={createPlanAction} className="grid gap-4 md:grid-cols-2">
            <FormField
              label="Plan name"
              hint="Use the show or stream name operators will recognize at a glance."
              className="md:col-span-2"
            >
              <Input name="name" placeholder="Saturday React Marathon" required />
            </FormField>
            <FormField
              label="Plan date"
              hint="The day this plan should appear in the planning queue."
            >
              <Input name="plan_date" type="date" required />
            </FormField>
            <FormField
              label="Timezone"
              hint="Keep this aligned with when the broadcaster expects clock-based segments to start."
            >
              <Input name="timezone" defaultValue="UTC" />
            </FormField>
            <FormField
              label="Planned start"
              hint="Used as the anchor for the stream day and for any clock-based planning."
              className="md:col-span-2"
            >
              <Input name="planned_start_at" type="datetime-local" />
            </FormField>
            <FormField
              label="Status"
              hint="Draft is editable planning state, scheduled is ready, live is currently running."
            >
              <Select name="status" defaultValue="draft">
                <option value="draft">Draft</option>
                <option value="scheduled">Scheduled</option>
                <option value="live">Live</option>
                <option value="completed">Completed</option>
                <option value="archived">Archived</option>
              </Select>
            </FormField>
            <label className="flex items-center gap-3 pt-8">
              <Switch name="auto_title_enabled" defaultChecked />
              <div>
                <span className="text-sm font-medium">Enable automatic titles</span>
                <p className="text-xs text-black/55 dark:text-white/55">
                  This becomes the plan default for segment-driven Twitch title updates.
                </p>
              </div>
            </label>
            <div className="md:col-span-2">
              <Button type="submit">Create plan</Button>
            </div>
          </form>
        </Card>
      </AppShell>
    );
  } catch {
    redirect("/login");
  }
}
