"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

import { TwitchLoginButton } from "@/components/twitch-login-button";
import { Card } from "@/components/ui/card";
import { authClient } from "@/lib/auth-client";

export default function LoginPage() {
  const router = useRouter();
  const { data: session, isPending } = authClient.useSession();

  useEffect(() => {
    if (session) {
      router.replace("/");
    }
  }, [router, session]);

  return (
    <div className="mx-auto flex min-h-screen max-w-lg items-center px-4">
      <Card className="w-full space-y-6 bg-white/95">
        <div className="space-y-2">
          <p className="text-xs uppercase tracking-[0.2em] text-black/45">
            Cuepoint
          </p>
          <h1 className="text-3xl font-semibold">Run the stream from one desk</h1>
          <p className="text-sm text-black/65">
            Sign in with Twitch to manage plans, title automation, overlays, and bot controls.
          </p>
        </div>
        <TwitchLoginButton />
        {isPending ? (
          <p className="text-xs text-black/45">Checking your session...</p>
        ) : null}
      </Card>
    </div>
  );
}
