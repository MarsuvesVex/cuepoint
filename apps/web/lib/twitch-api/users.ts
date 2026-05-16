import { helixRequest } from "@/lib/twitch-api/client";
import type { TwitchHelixResponse, TwitchUserProfile } from "@/lib/twitch-api/types";

export async function getCurrentTwitchUser(userId: string) {
  const response = await helixRequest<TwitchHelixResponse<TwitchUserProfile>>("/users", {
    userId,
  });

  return response.data[0] ?? null;
}
