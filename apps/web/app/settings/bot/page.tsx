import { redirect } from "next/navigation";

import { AppShell } from "@/components/app-shell";
import { FormField } from "@/components/form-field";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { requireSession } from "@/lib/auth";
import { saveBotSettingsAction } from "@/lib/actions";
import { db } from "@/lib/db";

export const dynamic = "force-dynamic";

export default async function BotSettingsPage() {
  try {
    const session = await requireSession();
    const result = await db.query(
      `SELECT auto_titles_default_enabled, react_commands_enabled, marker_commands_enabled,
              segment_advance_enabled, default_title_format
       FROM bot_settings
       WHERE user_id = $1`,
      [session.user.id],
    );
    const settings = result.rows[0] ?? {
      auto_titles_default_enabled: true,
      react_commands_enabled: true,
      marker_commands_enabled: true,
      segment_advance_enabled: true,
      default_title_format: "{original_title} | {segment_title}",
    };

    return (
      <AppShell>
        <Card className="max-w-3xl space-y-6">
          <div>
            <p className="text-xs uppercase tracking-[0.2em] text-black/45 dark:text-white/45">
              Bot settings
            </p>
            <h2 className="text-2xl font-semibold">Per-user automation defaults</h2>
            <p className="mt-2 text-sm text-black/60 dark:text-white/60">
              These toggles shape the runtime defaults the bot should start with before live session overrides change them.
            </p>
          </div>
          <form action={saveBotSettingsAction} className="space-y-4">
            <label className="flex items-center gap-3 rounded-2xl border border-black/10 bg-sand/50 px-4 py-4 dark:border-white/10 dark:bg-white/5">
              <Switch
                name="auto_titles_default_enabled"
                defaultChecked={settings.auto_titles_default_enabled}
              />
              <div>
                <span className="text-sm">Automatic title updates enabled by default</span>
                <p className="text-xs text-black/55 dark:text-white/55">
                  New live sessions start with title automation on unless disabled at the session level.
                </p>
              </div>
            </label>
            <label className="flex items-center gap-3 rounded-2xl border border-black/10 bg-sand/50 px-4 py-4 dark:border-white/10 dark:bg-white/5">
              <Switch
                name="react_commands_enabled"
                defaultChecked={settings.react_commands_enabled}
              />
              <div>
                <span className="text-sm">React commands enabled</span>
                <p className="text-xs text-black/55 dark:text-white/55">
                  Allows react-specific commands such as `!react` and `!watching` to operate by default.
                </p>
              </div>
            </label>
            <label className="flex items-center gap-3 rounded-2xl border border-black/10 bg-sand/50 px-4 py-4 dark:border-white/10 dark:bg-white/5">
              <Switch
                name="marker_commands_enabled"
                defaultChecked={settings.marker_commands_enabled}
              />
              <div>
                <span className="text-sm">Marker commands enabled</span>
                <p className="text-xs text-black/55 dark:text-white/55">
                  Enables timeline marker commands that later feed VOD clipping and FFmpeg output.
                </p>
              </div>
            </label>
            <label className="flex items-center gap-3 rounded-2xl border border-black/10 bg-sand/50 px-4 py-4 dark:border-white/10 dark:bg-white/5">
              <Switch
                name="segment_advance_enabled"
                defaultChecked={settings.segment_advance_enabled}
              />
              <div>
                <span className="text-sm">Segment advance commands enabled</span>
                <p className="text-xs text-black/55 dark:text-white/55">
                  Lets manual advancement commands move the active segment when automation should not wait for elapsed timing.
                </p>
              </div>
            </label>
            <FormField
              label="Default title format"
              hint="Allowed placeholders are `{original_title}` and `{segment_title}`. This is the fallback format for live sessions."
            >
              <Input
                name="default_title_format"
                defaultValue={settings.default_title_format}
              />
            </FormField>
            <Button type="submit">Save settings</Button>
          </form>
        </Card>
      </AppShell>
    );
  } catch {
    redirect("/login");
  }
}
