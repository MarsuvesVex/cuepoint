import { helixRequest } from "@/lib/twitch-api/client";
import type { TwitchHelixResponse, TwitchVideo } from "@/lib/twitch-api/types";

type GetVideosInput = {
  userId?: string;
  ids?: string[];
  channelUserId?: string;
  first?: number;
  after?: string;
};

// Intended future callers:
// - VOD cut/marker review pages
// - marker lookup flows that need to resolve recent broadcaster videos
export async function getVideos(input: GetVideosInput) {
  const query = new URLSearchParams();
  for (const id of input.ids ?? []) {
    query.append("id", id);
  }
  if (input.channelUserId) {
    query.set("user_id", input.channelUserId);
  }
  if (input.first) query.set("first", String(input.first));
  if (input.after) query.set("after", input.after);

  return helixRequest<TwitchHelixResponse<TwitchVideo>>("/videos", {
    userId: input.userId,
    appAccess: !input.userId,
    query,
  });
}
