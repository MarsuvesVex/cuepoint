import { Pool } from "pg";

declare global {
  // eslint-disable-next-line no-var
  var __cuepointWebPool: Pool | undefined;
  // eslint-disable-next-line no-var
  var __cuepointWebSchemaInit: Promise<void> | undefined;
}

export const db =
  global.__cuepointWebPool ??
  new Pool({
    connectionString:
      process.env.DATABASE_URL ??
      "postgres://cuepoint:cuepoint@127.0.0.1:5439/cuepoint?sslmode=disable",
  });

if (process.env.NODE_ENV !== "production") {
  global.__cuepointWebPool = db;
}

const bootstrapSQL = `
CREATE TABLE IF NOT EXISTS "user" (
  "id" TEXT PRIMARY KEY,
  "name" TEXT NOT NULL,
  "email" TEXT NOT NULL UNIQUE,
  "emailVerified" BOOLEAN NOT NULL DEFAULT FALSE,
  "image" TEXT,
  "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "session" (
  "id" TEXT PRIMARY KEY,
  "expiresAt" TIMESTAMPTZ NOT NULL,
  "token" TEXT NOT NULL UNIQUE,
  "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "ipAddress" TEXT,
  "userAgent" TEXT,
  "userId" TEXT NOT NULL REFERENCES "user"("id") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS "session_userId_idx" ON "session" ("userId");

CREATE TABLE IF NOT EXISTS "account" (
  "id" TEXT PRIMARY KEY,
  "accountId" TEXT NOT NULL,
  "providerId" TEXT NOT NULL,
  "userId" TEXT NOT NULL REFERENCES "user"("id") ON DELETE CASCADE,
  "accessToken" TEXT,
  "refreshToken" TEXT,
  "idToken" TEXT,
  "accessTokenExpiresAt" TIMESTAMPTZ,
  "refreshTokenExpiresAt" TIMESTAMPTZ,
  "scope" TEXT,
  "password" TEXT,
  "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "account_userId_idx" ON "account" ("userId");
CREATE UNIQUE INDEX IF NOT EXISTS "account_provider_account_idx" ON "account" ("providerId", "accountId");

CREATE TABLE IF NOT EXISTS "verification" (
  "id" TEXT PRIMARY KEY,
  "identifier" TEXT NOT NULL,
  "value" TEXT NOT NULL,
  "expiresAt" TIMESTAMPTZ NOT NULL,
  "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "verification_identifier_idx" ON "verification" ("identifier");

CREATE TABLE IF NOT EXISTS operator_profiles (
  user_id TEXT PRIMARY KEY,
  timezone TEXT NOT NULL DEFAULT 'UTC',
  twitch_broadcaster_id TEXT NOT NULL DEFAULT '',
  twitch_login TEXT NOT NULL DEFAULT '',
  twitch_display_name TEXT NOT NULL DEFAULT '',
  overlay_public_id TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS operator_profiles_twitch_login_idx
  ON operator_profiles (LOWER(twitch_login))
  WHERE twitch_login <> '';

CREATE UNIQUE INDEX IF NOT EXISTS operator_profiles_overlay_public_id_idx
  ON operator_profiles (overlay_public_id)
  WHERE overlay_public_id <> '';

CREATE TABLE IF NOT EXISTS bot_settings (
  user_id TEXT PRIMARY KEY,
  auto_titles_default_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  react_commands_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  marker_commands_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  segment_advance_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  default_title_format TEXT NOT NULL DEFAULT '{original_title} | {segment_title}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS stream_plans (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  plan_date DATE NOT NULL,
  name TEXT NOT NULL,
  timezone TEXT NOT NULL DEFAULT 'UTC',
  planned_start_at TIMESTAMPTZ,
  status TEXT NOT NULL,
  auto_title_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS stream_plans_user_date_idx
  ON stream_plans (user_id, plan_date DESC);

CREATE TABLE IF NOT EXISTS stream_plan_segments (
  id TEXT PRIMARY KEY,
  plan_id TEXT NOT NULL REFERENCES stream_plans(id) ON DELETE CASCADE,
  position INTEGER NOT NULL,
  segment_type TEXT NOT NULL,
  segment_title TEXT NOT NULL,
  notes TEXT NOT NULL DEFAULT '',
  stream_title_override TEXT NOT NULL DEFAULT '',
  timing_mode TEXT NOT NULL,
  start_at_local TIMESTAMPTZ,
  elapsed_offset_seconds INTEGER,
  duration_estimate_seconds INTEGER NOT NULL DEFAULT 0,
  youtube_url TEXT NOT NULL DEFAULT '',
  youtube_video_id TEXT NOT NULL DEFAULT '',
  youtube_video_title TEXT NOT NULL DEFAULT '',
  youtube_creator_name TEXT NOT NULL DEFAULT '',
  youtube_canonical_url TEXT NOT NULL DEFAULT '',
  metadata_status TEXT NOT NULL DEFAULT 'pending',
  metadata_error TEXT NOT NULL DEFAULT '',
  react_start_command_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  react_title_command_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  react_watching_command_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS stream_plan_segments_plan_position_idx
  ON stream_plan_segments (plan_id, position ASC);

CREATE TABLE IF NOT EXISTS stream_sessions (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  plan_id TEXT NOT NULL REFERENCES stream_plans(id) ON DELETE CASCADE,
  twitch_stream_id TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL,
  started_at TIMESTAMPTZ NOT NULL,
  ended_at TIMESTAMPTZ,
  original_title TEXT NOT NULL DEFAULT '',
  current_title TEXT NOT NULL DEFAULT '',
  current_segment_id TEXT NOT NULL DEFAULT '',
  current_segment_started_at TIMESTAMPTZ,
  auto_title_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  title_format_override TEXT NOT NULL DEFAULT '',
  last_synced_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS stream_sessions_user_status_idx
  ON stream_sessions (user_id, status);

CREATE TABLE IF NOT EXISTS timeline_markers (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  plan_id TEXT NOT NULL REFERENCES stream_plans(id) ON DELETE CASCADE,
  session_id TEXT NOT NULL REFERENCES stream_sessions(id) ON DELETE CASCADE,
  segment_id TEXT NOT NULL DEFAULT '',
  kind TEXT NOT NULL,
  label TEXT NOT NULL,
  paired_marker_id TEXT NOT NULL DEFAULT '',
  stream_elapsed_seconds INTEGER NOT NULL DEFAULT 0,
  source TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS timeline_markers_session_created_idx
  ON timeline_markers (session_id, created_at DESC);

CREATE TABLE IF NOT EXISTS automation_jobs (
  id TEXT PRIMARY KEY,
  job_type TEXT NOT NULL,
  resource_type TEXT NOT NULL,
  resource_id TEXT NOT NULL,
  payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  status TEXT NOT NULL,
  error TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS automation_jobs_status_idx
  ON automation_jobs (status, created_at ASC);

CREATE TABLE IF NOT EXISTS dev_user_settings (
  user_id TEXT PRIMARY KEY,
  force_live_state BOOLEAN NOT NULL DEFAULT FALSE,
  forced_stream_title TEXT NOT NULL DEFAULT '',
  forced_stream_id TEXT NOT NULL DEFAULT '',
  forced_started_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`;

export async function ensureWebSchema() {
  if (!global.__cuepointWebSchemaInit) {
    global.__cuepointWebSchemaInit = db.query(bootstrapSQL).then(() => undefined);
  }

  return global.__cuepointWebSchemaInit;
}
