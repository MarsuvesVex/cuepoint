"use client";

import { useTransition } from "react";

import { Button } from "@/components/ui/button";
import { authClient } from "@/lib/auth-client";

export function TwitchLoginButton() {
  const [pending, startTransition] = useTransition();

  return (
    <Button
      onClick={() =>
        startTransition(async () => {
          await authClient.signIn.social({
            provider: "twitch",
            callbackURL: "/",
          });
        })
      }
      disabled={pending}
      className="w-full"
    >
      {pending ? "Redirecting..." : "Continue with Twitch"}
    </Button>
  );
}
