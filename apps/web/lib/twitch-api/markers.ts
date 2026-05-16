import { helixRequest, requireBroadcasterProfile } from "@/lib/twitch-api/client";
import type { TwitchHelixResponse, TwitchStreamMarker } from "@/lib/twitch-api/types";

type MarkerResponse = {
  data: Array<TwitchStreamMarker>;
};

type MarkerListingResponse = {
  data: Array<{
    user_id: string;
    user_name: string;
    user_login: string;
    videos: Array<{
      video_id: string;
      markers: TwitchStreamMarker[];
    }>;
  }>;
  pagination?: {
    cursor?: string;
  };
};

// Intended future callers:
// - bot/runtime marker mirroring in apps/api/internal/api/runtime.go
// - VOD tooling pages once marker-to-video workflows are exposed in the UI
export async function createStreamMarker(userId: string, description?: string) {
  const profile = await requireBroadcasterProfile(userId);
  const response = await helixRequest<MarkerResponse>("/streams/markers", {
    method: "POST",
    userId,
    body: {
      user_id: profile.twitch_broadcaster_id,
      ...(description ? { description } : {}),
    },
  });

  return response.data[0] ?? null;
}

export async function getStreamMarkers(input: {
  userId: string;
  videoId?: string;
  first?: number;
  after?: string;
  before?: string;
}) {
  const profile = await requireBroadcasterProfile(input.userId);
  const query = new URLSearchParams();
  if (input.videoId) {
    query.set("video_id", input.videoId);
  } else {
    query.set("user_id", profile.twitch_broadcaster_id);
  }
  if (input.first) query.set("first", String(input.first));
  if (input.after) query.set("after", input.after);
  if (input.before) query.set("before", input.before);

  return helixRequest<MarkerListingResponse>("/streams/markers", {
    userId: input.userId,
    query,
  });
}
