import { redirect } from "next/navigation";

import { AppShell } from "@/components/app-shell";
import { FormField } from "@/components/form-field";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { saveDevSettingsAction } from "@/lib/actions";
import { requireSession } from "@/lib/auth";
import { getDevSettings } from "@/lib/data";

export const dynamic = "force-dynamic";

export default async function DevSettingsPage() {
  try {
    const session = await requireSession();
    const settings = await getDevSettings(session.user.id);

    return (
      <AppShell>
        <Card className="max-w-3xl space-y-6">
          <div>
            <p className="text-xs uppercase tracking-[0.2em] text-black/45 dark:text-white/45">
              Dev settings
            </p>
            <h2 className="text-2xl font-semibold">Local runtime overrides</h2>
            <p className="mt-2 text-sm text-black/60 dark:text-white/60">
              These values are for testing the app and the internal Twitch bridge without waiting for a real live broadcast.
            </p>
          </div>

          <form action={saveDevSettingsAction} className="space-y-5">
            <label className="flex items-center gap-3 rounded-2xl border border-black/10 bg-sand/50 px-4 py-4 dark:border-white/10 dark:bg-white/5">
              <Switch
                name="force_live_state"
                defaultChecked={settings.force_live_state}
              />
              <div>
                <p className="text-sm font-medium">Force live state</p>
                <p className="text-xs text-black/55 dark:text-white/55">
                  When enabled, the Twitch stream lookup returns a synthetic live stream for this user.
                </p>
              </div>
            </label>

            <div className="grid gap-4 md:grid-cols-2">
              <FormField
                label="Forced stream title"
                hint="Used as the live stream title returned to the runtime when forced live state is on."
                className="md:col-span-2"
              >
                <Input
                  name="forced_stream_title"
                  defaultValue={settings.forced_stream_title}
                  placeholder="Dev stream title"
                />
              </FormField>

              <FormField
                label="Forced stream ID"
                hint="A stable fake Twitch stream ID for bot/runtime testing."
              >
                <Input
                  name="forced_stream_id"
                  defaultValue={settings.forced_stream_id}
                  placeholder="dev-stream-1"
                />
              </FormField>

              <FormField
                label="Forced started at"
                hint="Use this to test elapsed segment activation without waiting in real time."
              >
                <Input
                  name="forced_started_at"
                  type="datetime-local"
                  defaultValue={
                    settings.forced_started_at
                      ? new Date(settings.forced_started_at).toISOString().slice(0, 16)
                      : ""
                  }
                />
              </FormField>
            </div>

            <Button type="submit">Save dev settings</Button>
          </form>
        </Card>
      </AppShell>
    );
  } catch {
    redirect("/login");
  }
}
