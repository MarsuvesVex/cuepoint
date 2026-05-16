import { betterAuth } from "better-auth";
import { toNextJsHandler } from "better-auth/next-js";
import { headers } from "next/headers";

import { db, ensureWebSchema } from "@/lib/db";

export const auth = betterAuth({
  appName: "Cuepoint Control Center",
  baseURL: process.env.BETTER_AUTH_URL ?? process.env.WEB_BASE_URL ?? "http://localhost:3000",
  database: db,
  advanced: {
    cookiePrefix: "cuepoint",
  },
  account: {
    encryptOAuthTokens: true,
    updateAccountOnSignIn: true,
  },
  socialProviders: {
    twitch: {
      clientId: process.env.TWITCH_CLIENT_ID as string,
      clientSecret: process.env.TWITCH_CLIENT_SECRET as string,
      scope: [
        "user:read:email",
        "channel:manage:broadcast",
      ],
    },
  },
});

export const authHandler = toNextJsHandler(auth);

export async function requireSession() {
  await ensureWebSchema();
  const session = await auth.api.getSession({
    headers: await headers(),
  });

  if (!session) {
    throw new Error("UNAUTHENTICATED");
  }

  return session;
}
