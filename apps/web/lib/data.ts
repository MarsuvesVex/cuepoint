import { db } from "@/lib/db";

export async function listPlans(userId: string) {
  const result = await db.query(
    `SELECT id, name, plan_date, timezone, planned_start_at, status, auto_title_enabled
     FROM stream_plans
     WHERE user_id = $1
     ORDER BY plan_date DESC, created_at DESC`,
    [userId],
  );

  return result.rows;
}

export async function getPlan(planId: string, userId: string) {
  const plan = await db.query(
    `SELECT id, name, plan_date, timezone, planned_start_at, status, auto_title_enabled
     FROM stream_plans
     WHERE id = $1 AND user_id = $2`,
    [planId, userId],
  );
  const segments = await db.query(
    `SELECT id, position, segment_type, segment_title, notes, stream_title_override, timing_mode,
            start_at_local, elapsed_offset_seconds, duration_estimate_seconds, youtube_url,
            youtube_video_title, youtube_creator_name, youtube_canonical_url, metadata_status,
            react_start_command_enabled, react_title_command_enabled, react_watching_command_enabled
     FROM stream_plan_segments
     WHERE plan_id = $1
     ORDER BY position ASC, created_at ASC`,
    [planId],
  );

  return {
    plan: plan.rows[0] ?? null,
    segments: segments.rows,
  };
}

export async function getDashboardData(userId: string) {
  const profileResult = await db.query(
    `SELECT twitch_login, twitch_display_name, overlay_public_id
     FROM operator_profiles
     WHERE user_id = $1`,
    [userId],
  );
  const settingsResult = await db.query(
    `SELECT auto_titles_default_enabled, react_commands_enabled, marker_commands_enabled,
            segment_advance_enabled, default_title_format
     FROM bot_settings
     WHERE user_id = $1`,
    [userId],
  );
  const plans = await listPlans(userId);
  const liveSession = await db.query(
    `SELECT id, current_title, original_title, current_segment_id, auto_title_enabled
     FROM stream_sessions
     WHERE user_id = $1 AND status = 'live'
     ORDER BY started_at DESC
     LIMIT 1`,
    [userId],
  );
  const devSettings = await db.query(
    `SELECT force_live_state, forced_stream_title, forced_stream_id, forced_started_at
     FROM dev_user_settings
     WHERE user_id = $1`,
    [userId],
  );

  return {
    profile: profileResult.rows[0] ?? null,
    settings: settingsResult.rows[0] ?? null,
    plans: plans.slice(0, 5),
    liveSession: liveSession.rows[0] ?? null,
    devSettings: devSettings.rows[0] ?? null,
  };
}

export async function getOverlayData(overlayPublicID: string) {
  const result = await db.query(
    `SELECT
       s.current_segment_id,
       seg.segment_title,
       seg.segment_type,
       seg.youtube_video_title,
       seg.youtube_creator_name,
       seg.youtube_canonical_url
     FROM operator_profiles p
     JOIN stream_sessions s ON s.user_id = p.user_id AND s.status = 'live'
     LEFT JOIN stream_plan_segments seg ON seg.id = s.current_segment_id
     WHERE p.overlay_public_id = $1
     ORDER BY s.started_at DESC
     LIMIT 1`,
    [overlayPublicID],
  );

  return result.rows[0] ?? null;
}

export async function getDevSettings(userId: string) {
  const result = await db.query(
    `SELECT force_live_state, forced_stream_title, forced_stream_id, forced_started_at
     FROM dev_user_settings
     WHERE user_id = $1`,
    [userId],
  );

  return (
    result.rows[0] ?? {
      force_live_state: false,
      forced_stream_title: "",
      forced_stream_id: "",
      forced_started_at: null,
    }
  );
}
