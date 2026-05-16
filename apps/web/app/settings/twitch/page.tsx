import { redirect } from "next/navigation";

import { AppShell } from "@/components/app-shell";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { requireSession } from "@/lib/auth";
import { syncTwitchProfileAction } from "@/lib/actions";
import { getOperatorProfile } from "@/lib/twitch";

export const dynamic = "force-dynamic";

export default async function TwitchSettingsPage() {
  try {
    const session = await requireSession();
    const profile = await getOperatorProfile(session.user.id);

    return (
      <AppShell>
        <Card className="max-w-3xl space-y-6">
          <div>
            <p className="text-xs uppercase tracking-[0.2em] text-black/45 dark:text-white/45">
              Twitch
            </p>
            <h2 className="text-2xl font-semibold">Account link and broadcaster identity</h2>
            <p className="mt-2 text-sm text-black/60 dark:text-white/60">
              This page owns the broadcaster identity that later Helix calls, overlays, markers, and schedule sync will use.
            </p>
          </div>
          <div className="rounded-2xl bg-sand/60 p-4 text-sm text-black/70 dark:bg-white/5 dark:text-white/70">
            {profile ? (
              <div className="space-y-2">
                <p>Login: {profile.twitch_login}</p>
                <p>Display name: {profile.twitch_display_name}</p>
                <p>Broadcaster ID: {profile.twitch_broadcaster_id}</p>
                <p>Overlay public ID: {profile.overlay_public_id}</p>
              </div>
            ) : (
              <p>No broadcaster profile synced yet. Sign in with Twitch first, then sync your broadcaster profile.</p>
            )}
          </div>
          <form action={syncTwitchProfileAction}>
            <Button type="submit">Sync Twitch profile</Button>
          </form>
        </Card>
      </AppShell>
    );
  } catch {
    redirect("/login");
  }
}
