import { auth } from "@/lib/auth";
import { db } from "@/lib/db";

type TwitchAccessToken = {
  accessToken: string;
};

export async function getTwitchAccessToken(userId: string): Promise<string> {
  const token = (await auth.api.getAccessToken({
    body: {
      providerId: "twitch",
      userId,
    },
  })) as TwitchAccessToken | null;

  if (!token?.accessToken) {
    throw new Error("No Twitch access token found");
  }

  return token.accessToken;
}

export async function getTwitchAppAccessToken(): Promise<string> {
  const clientId = process.env.TWITCH_CLIENT_ID ?? "";
  const clientSecret = process.env.TWITCH_CLIENT_SECRET ?? "";
  const response = await fetch("https://id.twitch.tv/oauth2/token", {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: new URLSearchParams({
      client_id: clientId,
      client_secret: clientSecret,
      grant_type: "client_credentials",
    }),
    cache: "no-store",
  });

  if (!response.ok) {
    throw new Error(`Twitch app token fetch failed: ${response.status}`);
  }

  const data = (await response.json()) as { access_token?: string };
  if (!data.access_token) {
    throw new Error("Twitch app token missing");
  }
  return data.access_token;
}

export async function getOperatorProfile(userId: string) {
  const result = await db.query(
    `SELECT user_id, timezone, twitch_broadcaster_id, twitch_login, twitch_display_name, overlay_public_id
     FROM operator_profiles
     WHERE user_id = $1`,
    [userId],
  );

  return result.rows[0] ?? null;
}
