"use server";

import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";

import { requireSession } from "@/lib/auth";
import { db } from "@/lib/db";
import { syncTwitchProfile } from "@/lib/twitch";

export async function createPlanAction(formData: FormData) {
  const session = await requireSession();
  const id = crypto.randomUUID();
  await db.query(
    `INSERT INTO stream_plans (
       id, user_id, plan_date, name, timezone, planned_start_at, status, auto_title_enabled, created_at, updated_at
     )
     VALUES ($1, $2, $3, $4, $5, NULLIF($6, '')::timestamptz, $7, $8, NOW(), NOW())`,
    [
      id,
      session.user.id,
      formData.get("plan_date"),
      formData.get("name"),
      formData.get("timezone") || "UTC",
      formData.get("planned_start_at"),
      formData.get("status") || "draft",
      formData.get("auto_title_enabled") === "on",
    ],
  );

  revalidatePath("/");
  redirect(`/plans/${id}`);
}

export async function savePlanAction(formData: FormData) {
  const session = await requireSession();
  const planID = String(formData.get("plan_id"));
  const existing = await db.query(
    `SELECT plan_date
     FROM stream_plans
     WHERE id = $1 AND user_id = $2`,
    [planID, session.user.id],
  );
  const currentPlanDate = existing.rows[0]?.plan_date;
  const incomingPlanDate = String(formData.get("plan_date") || "").trim();
  const planDate = incomingPlanDate || currentPlanDate;

  if (!planDate) {
    throw new Error("plan date is required");
  }

  await db.query(
     `UPDATE stream_plans
      SET name = $3,
          plan_date = $4,
          timezone = $5,
         planned_start_at = NULLIF($6, '')::timestamptz,
         status = $7,
         auto_title_enabled = $8,
         updated_at = NOW()
     WHERE id = $1 AND user_id = $2`,
    [
      planID,
      session.user.id,
      formData.get("name"),
      planDate,
      formData.get("timezone") || "UTC",
      formData.get("planned_start_at"),
      formData.get("status") || "draft",
      formData.get("auto_title_enabled") === "on",
    ],
  );
  revalidatePath(`/plans/${formData.get("plan_id")}`);
  revalidatePath("/plans");
}

export async function addSegmentAction(formData: FormData) {
  const session = await requireSession();
  const planID = String(formData.get("plan_id"));
  const plan = await db.query(
    `SELECT user_id FROM stream_plans WHERE id = $1`,
    [planID],
  );
  if (plan.rows[0]?.user_id !== session.user.id) {
    throw new Error("Unauthorized");
  }

  const id = crypto.randomUUID();
  const positionResult = await db.query(
    `SELECT COALESCE(MAX(position), 0) + 1 AS next_position
     FROM stream_plan_segments
     WHERE plan_id = $1`,
    [planID],
  );
  const position = Number(positionResult.rows[0]?.next_position ?? 1);
  const timingMode = String(formData.get("timing_mode") || "elapsed");
  const youtubeURL = String(formData.get("youtube_url") || "");

  await db.query(
    `INSERT INTO stream_plan_segments (
       id, plan_id, position, segment_type, segment_title, notes, stream_title_override, timing_mode,
       start_at_local, elapsed_offset_seconds, duration_estimate_seconds, youtube_url, metadata_status,
       react_start_command_enabled, react_title_command_enabled, react_watching_command_enabled, created_at, updated_at
     )
     VALUES (
       $1, $2, $3, $4, $5, $6, $7, $8,
       NULLIF($9, '')::timestamptz, NULLIF($10, '')::int, $11, $12, $13,
       $14, $15, $16, NOW(), NOW()
     )`,
    [
      id,
      planID,
      position,
      formData.get("segment_type") || "standard",
      formData.get("segment_title"),
      formData.get("notes") || "",
      formData.get("stream_title_override") || "",
      timingMode,
      formData.get("start_at_local"),
      formData.get("elapsed_offset_seconds"),
      Number(formData.get("duration_estimate_seconds") || 0),
      youtubeURL,
      youtubeURL ? "pending" : "ready",
      formData.get("react_start_command_enabled") === "on",
      formData.get("react_title_command_enabled") === "on",
      formData.get("react_watching_command_enabled") === "on",
    ],
  );

  if (youtubeURL) {
    await db.query(
      `INSERT INTO automation_jobs (
         id, job_type, resource_type, resource_id, payload_json, status, error, created_at, updated_at
       )
       VALUES ($1, 'youtube_metadata_sync', 'stream_plan_segment', $2, '{}'::jsonb, 'pending', '', NOW(), NOW())`,
      [crypto.randomUUID(), id],
    );
  }

  revalidatePath(`/plans/${planID}`);
}

export async function saveBotSettingsAction(formData: FormData) {
  const session = await requireSession();
  await db.query(
    `INSERT INTO bot_settings (
       user_id, auto_titles_default_enabled, react_commands_enabled, marker_commands_enabled,
       segment_advance_enabled, default_title_format, created_at, updated_at
     )
     VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
     ON CONFLICT (user_id) DO UPDATE SET
       auto_titles_default_enabled = EXCLUDED.auto_titles_default_enabled,
       react_commands_enabled = EXCLUDED.react_commands_enabled,
       marker_commands_enabled = EXCLUDED.marker_commands_enabled,
       segment_advance_enabled = EXCLUDED.segment_advance_enabled,
       default_title_format = EXCLUDED.default_title_format,
       updated_at = NOW()`,
    [
      session.user.id,
      formData.get("auto_titles_default_enabled") === "on",
      formData.get("react_commands_enabled") === "on",
      formData.get("marker_commands_enabled") === "on",
      formData.get("segment_advance_enabled") === "on",
      formData.get("default_title_format") || "{original_title} | {segment_title}",
    ],
  );

  revalidatePath("/settings/bot");
  revalidatePath("/");
}

export async function syncTwitchProfileAction() {
  const session = await requireSession();
  await syncTwitchProfile(session.user.id);
  revalidatePath("/settings/twitch");
  revalidatePath("/");
}

export async function saveDevSettingsAction(formData: FormData) {
  const session = await requireSession();
  await db.query(
    `INSERT INTO dev_user_settings (
       user_id, force_live_state, forced_stream_title, forced_stream_id, forced_started_at, created_at, updated_at
     )
     VALUES ($1, $2, $3, $4, NULLIF($5, '')::timestamptz, NOW(), NOW())
     ON CONFLICT (user_id) DO UPDATE SET
       force_live_state = EXCLUDED.force_live_state,
       forced_stream_title = EXCLUDED.forced_stream_title,
       forced_stream_id = EXCLUDED.forced_stream_id,
       forced_started_at = EXCLUDED.forced_started_at,
       updated_at = NOW()`,
    [
      session.user.id,
      formData.get("force_live_state") === "on",
      formData.get("forced_stream_title") || "",
      formData.get("forced_stream_id") || "",
      formData.get("forced_started_at") || "",
    ],
  );

  revalidatePath("/settings/dev");
  revalidatePath("/");
}
