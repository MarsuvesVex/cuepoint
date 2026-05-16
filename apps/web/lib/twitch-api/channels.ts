import { helixRequest, requireBroadcasterProfile } from "@/lib/twitch-api/client";
import type { TwitchChannel, TwitchHelixResponse } from "@/lib/twitch-api/types";

export async function getChannelInformation(userId: string) {
  const profile = await requireBroadcasterProfile(userId);
  const response = await helixRequest<TwitchHelixResponse<TwitchChannel>>("/channels", {
    userId,
    query: new URLSearchParams({
      broadcaster_id: profile.twitch_broadcaster_id,
    }),
  });

  return response.data[0] ?? null;
}

export async function updateChannelTitle(userId: string, title: string) {
  const profile = await requireBroadcasterProfile(userId);
  await helixRequest("/channels", {
    method: "PATCH",
    userId,
    query: new URLSearchParams({
      broadcaster_id: profile.twitch_broadcaster_id,
    }),
    body: { title },
  });
}
