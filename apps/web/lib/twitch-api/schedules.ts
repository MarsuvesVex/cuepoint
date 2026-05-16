import { helixRequest, requireBroadcasterProfile } from "@/lib/twitch-api/client";
import type { TwitchSchedule } from "@/lib/twitch-api/types";

type ScheduleResponse = {
  data: TwitchSchedule;
  pagination?: {
    cursor?: string;
  };
};

type UpsertSegmentInput = {
  userId: string;
  segmentId?: string;
  startTime: string;
  timezone: string;
  durationMinutes: number;
  title?: string;
  isRecurring?: boolean;
  categoryId?: string;
};

// Intended future callers:
// - plan editor sync/export actions
// - Twitch schedule reconciliation jobs once plans map onto broadcaster schedule segments
export async function getChannelSchedule(input: {
  userId: string;
  startTime?: string;
  ids?: string[];
  first?: number;
  after?: string;
}) {
  const profile = await requireBroadcasterProfile(input.userId);
  const query = new URLSearchParams({
    broadcaster_id: profile.twitch_broadcaster_id,
  });
  if (input.startTime) query.set("start_time", input.startTime);
  if (input.first) query.set("first", String(input.first));
  if (input.after) query.set("after", input.after);
  for (const id of input.ids ?? []) {
    query.append("id", id);
  }

  return helixRequest<ScheduleResponse>("/schedule", {
    userId: input.userId,
    query,
  });
}

export async function createScheduleSegment(input: UpsertSegmentInput) {
  const profile = await requireBroadcasterProfile(input.userId);
  return helixRequest("/schedule/segment", {
    method: "POST",
    userId: input.userId,
    query: new URLSearchParams({
      broadcaster_id: profile.twitch_broadcaster_id,
    }),
    body: {
      start_time: input.startTime,
      timezone: input.timezone,
      duration: input.durationMinutes,
      ...(input.title ? { title: input.title } : {}),
      ...(input.isRecurring !== undefined ? { is_recurring: input.isRecurring } : {}),
      ...(input.categoryId ? { category_id: input.categoryId } : {}),
    },
  });
}

export async function updateScheduleSegment(input: UpsertSegmentInput) {
  if (!input.segmentId) {
    throw new Error("segmentId is required");
  }
  const profile = await requireBroadcasterProfile(input.userId);
  return helixRequest("/schedule/segment", {
    method: "PATCH",
    userId: input.userId,
    query: new URLSearchParams({
      broadcaster_id: profile.twitch_broadcaster_id,
      id: input.segmentId,
    }),
    body: {
      start_time: input.startTime,
      timezone: input.timezone,
      duration: input.durationMinutes,
      ...(input.title ? { title: input.title } : {}),
      ...(input.categoryId ? { category_id: input.categoryId } : {}),
      ...(input.isRecurring !== undefined ? { is_recurring: input.isRecurring } : {}),
    },
  });
}

export async function deleteScheduleSegment(userId: string, segmentId: string) {
  const profile = await requireBroadcasterProfile(userId);
  return helixRequest("/schedule/segment", {
    method: "DELETE",
    userId,
    query: new URLSearchParams({
      broadcaster_id: profile.twitch_broadcaster_id,
      id: segmentId,
    }),
  });
}

export async function updateScheduleVacation(input: {
  userId: string;
  timezone: string;
  startTime: string;
  endTime: string;
  enabled: boolean;
}) {
  const profile = await requireBroadcasterProfile(input.userId);
  return helixRequest("/schedule/settings", {
    method: "PATCH",
    userId: input.userId,
    query: new URLSearchParams({
      broadcaster_id: profile.twitch_broadcaster_id,
    }),
    body: {
      is_vacation_enabled: input.enabled,
      vacation_start_time: input.startTime,
      vacation_end_time: input.endTime,
      timezone: input.timezone,
    },
  });
}
