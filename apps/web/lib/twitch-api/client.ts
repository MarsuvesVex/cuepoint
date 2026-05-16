import {
  getOperatorProfile,
  getTwitchAccessToken,
  getTwitchAppAccessToken,
} from "@/lib/twitch-auth";

const HELIX_BASE_URL = "https://api.twitch.tv/helix";

type RequestOptions = {
  method?: "GET" | "POST" | "PATCH" | "DELETE";
  query?: URLSearchParams;
  body?: unknown;
  userId?: string;
  appAccess?: boolean;
};

export async function helixRequest<T>(
  path: string,
  { method = "GET", query, body, userId, appAccess = false }: RequestOptions = {},
): Promise<T> {
  const accessToken = appAccess
    ? await getTwitchAppAccessToken()
    : await getTwitchAccessToken(requiredUserId(userId));

  const url = new URL(`${HELIX_BASE_URL}${path}`);
  if (query) {
    query.forEach((value, key) => url.searchParams.append(key, value));
  }

  const response = await fetch(url, {
    method,
    headers: {
      Authorization: `Bearer ${accessToken}`,
      "Client-Id": process.env.TWITCH_CLIENT_ID ?? "",
      ...(body ? { "Content-Type": "application/json" } : {}),
    },
    body: body ? JSON.stringify(body) : undefined,
    cache: "no-store",
  });

  if (!response.ok) {
    throw new Error(`Twitch Helix request failed: ${response.status} ${path}`);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  const text = await response.text();
  if (!text.trim()) {
    return undefined as T;
  }

  return JSON.parse(text) as T;
}

export async function requireBroadcasterProfile(userId: string) {
  const profile = await getOperatorProfile(userId);
  if (!profile?.twitch_broadcaster_id) {
    throw new Error("Twitch broadcaster is not linked");
  }
  return profile;
}

function requiredUserId(userId?: string) {
  if (!userId) {
    throw new Error("userId is required for this Twitch request");
  }
  return userId;
}
