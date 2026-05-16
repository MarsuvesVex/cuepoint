import { db } from "@/lib/db";
import { getOperatorProfile, getTwitchAccessToken, getTwitchAppAccessToken } from "@/lib/twitch-auth";
import { getChannelInformation, updateChannelTitle as patchChannelTitle } from "@/lib/twitch-api/channels";
import { createStreamMarker, getStreamMarkers } from "@/lib/twitch-api/markers";
import {
  createScheduleSegment,
  deleteScheduleSegment,
  getChannelSchedule,
  updateScheduleSegment,
  updateScheduleVacation,
} from "@/lib/twitch-api/schedules";
import { getLiveStream } from "@/lib/twitch-api/streams";
import { getCurrentTwitchUser } from "@/lib/twitch-api/users";
import { getVideos } from "@/lib/twitch-api/videos";

export { getOperatorProfile, getTwitchAccessToken, getTwitchAppAccessToken };

export async function syncTwitchProfile(userId: string) {
  const user = await getCurrentTwitchUser(userId);
  if (!user) {
    throw new Error("Twitch profile not found");
  }

  const overlayPublicID = crypto.randomUUID();
  await db.query(
    `INSERT INTO operator_profiles (
       user_id, timezone, twitch_broadcaster_id, twitch_login, twitch_display_name, overlay_public_id, created_at, updated_at
     )
     VALUES ($1, 'UTC', $2, $3, $4, $5, NOW(), NOW())
     ON CONFLICT (user_id) DO UPDATE SET
       twitch_broadcaster_id = EXCLUDED.twitch_broadcaster_id,
       twitch_login = EXCLUDED.twitch_login,
       twitch_display_name = EXCLUDED.twitch_display_name,
       overlay_public_id = CASE
         WHEN operator_profiles.overlay_public_id = '' THEN EXCLUDED.overlay_public_id
         ELSE operator_profiles.overlay_public_id
       END,
       updated_at = NOW()`,
    [userId, user.id, user.login, user.display_name, overlayPublicID],
  );

  return getOperatorProfile(userId);
}

export async function getChannelState(userId: string) {
  const channel = await getChannelInformation(userId);
  return { title: channel?.title ?? "" };
}

export async function updateChannelTitle(userId: string, title: string) {
  await patchChannelTitle(userId, title);
}

export async function getStreamState(userId: string) {
  const override = await db.query(
    `SELECT force_live_state, forced_stream_title, forced_stream_id, forced_started_at
     FROM dev_user_settings
     WHERE user_id = $1`,
    [userId],
  );

  const devSettings = override.rows[0];
  if (devSettings?.force_live_state) {
    return {
      live: true,
      streamId: devSettings.forced_stream_id || "dev-stream-1",
      title: devSettings.forced_stream_title || "Forced live state",
      startedAt: devSettings.forced_started_at
        ? new Date(devSettings.forced_started_at).toISOString()
        : new Date().toISOString(),
    };
  }

  const streamData = await getLiveStream(userId);
  if (!streamData) {
    return {
      live: false,
      streamId: "",
      title: "",
      startedAt: new Date(0).toISOString(),
    };
  }

  return {
    live: true,
    streamId: streamData.id,
    title: streamData.title,
    startedAt: streamData.started_at,
  };
}

export const twitchMarkers = {
  create: createStreamMarker,
  list: getStreamMarkers,
};

export const twitchVideos = {
  list: getVideos,
};

export const twitchSchedule = {
  get: getChannelSchedule,
  createSegment: createScheduleSegment,
  updateSegment: updateScheduleSegment,
  deleteSegment: deleteScheduleSegment,
  updateVacation: updateScheduleVacation,
};
