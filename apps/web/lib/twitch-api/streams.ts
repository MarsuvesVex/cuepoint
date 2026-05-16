import { helixRequest, requireBroadcasterProfile } from "@/lib/twitch-api/client";
import type { TwitchHelixResponse, TwitchStream } from "@/lib/twitch-api/types";

export async function getLiveStream(userId: string) {
  const profile = await requireBroadcasterProfile(userId);
  const response = await helixRequest<TwitchHelixResponse<TwitchStream>>("/streams", {
    userId,
    query: new URLSearchParams({
      user_login: profile.twitch_login,
    }),
  });

  return response.data[0] ?? null;
}
