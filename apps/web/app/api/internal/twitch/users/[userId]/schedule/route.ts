import { NextRequest, NextResponse } from "next/server";

import { requireInternalService } from "@/lib/internal-auth";
import { twitchSchedule } from "@/lib/twitch";

export const dynamic = "force-dynamic";

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ userId: string }> },
) {
  try {
    await requireInternalService();
    const { userId } = await params;
    const startTime = request.nextUrl.searchParams.get("startTime") ?? undefined;
    const ids = request.nextUrl.searchParams.getAll("id");
    const first = request.nextUrl.searchParams.get("first");
    const after = request.nextUrl.searchParams.get("after") ?? undefined;

    return NextResponse.json(
      await twitchSchedule.get({
        userId,
        startTime,
        ids: ids.length > 0 ? ids : undefined,
        first: first ? Number(first) : undefined,
        after,
      }),
    );
  } catch (error) {
    return NextResponse.json(
      { error: error instanceof Error ? error.message : "Unknown error" },
      { status: 400 },
    );
  }
}

export async function PATCH(
  request: Request,
  { params }: { params: Promise<{ userId: string }> },
) {
  try {
    await requireInternalService();
    const { userId } = await params;
    const body = (await request.json()) as {
      timezone: string;
      startTime: string;
      endTime: string;
      enabled: boolean;
    };

    return NextResponse.json(
      await twitchSchedule.updateVacation({
        userId,
        timezone: body.timezone,
        startTime: body.startTime,
        endTime: body.endTime,
        enabled: body.enabled,
      }),
    );
  } catch (error) {
    return NextResponse.json(
      { error: error instanceof Error ? error.message : "Unknown error" },
      { status: 400 },
    );
  }
}
