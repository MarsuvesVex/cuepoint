import { NextResponse } from "next/server";

import { requireInternalService } from "@/lib/internal-auth";
import { twitchSchedule } from "@/lib/twitch";

export const dynamic = "force-dynamic";

export async function POST(
  request: Request,
  { params }: { params: Promise<{ userId: string }> },
) {
  try {
    await requireInternalService();
    const { userId } = await params;
    const body = (await request.json()) as {
      startTime: string;
      timezone: string;
      durationMinutes: number;
      title?: string;
      isRecurring?: boolean;
      categoryId?: string;
    };

    return NextResponse.json(
      await twitchSchedule.createSegment({
        userId,
        startTime: body.startTime,
        timezone: body.timezone,
        durationMinutes: body.durationMinutes,
        title: body.title,
        isRecurring: body.isRecurring,
        categoryId: body.categoryId,
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
      segmentId: string;
      startTime: string;
      timezone: string;
      durationMinutes: number;
      title?: string;
      isRecurring?: boolean;
      categoryId?: string;
    };

    return NextResponse.json(
      await twitchSchedule.updateSegment({
        userId,
        segmentId: body.segmentId,
        startTime: body.startTime,
        timezone: body.timezone,
        durationMinutes: body.durationMinutes,
        title: body.title,
        isRecurring: body.isRecurring,
        categoryId: body.categoryId,
      }),
    );
  } catch (error) {
    return NextResponse.json(
      { error: error instanceof Error ? error.message : "Unknown error" },
      { status: 400 },
    );
  }
}

export async function DELETE(
  request: Request,
  { params }: { params: Promise<{ userId: string }> },
) {
  try {
    await requireInternalService();
    const { userId } = await params;
    const url = new URL(request.url);
    const segmentId = url.searchParams.get("segmentId");
    if (!segmentId) {
      return NextResponse.json({ error: "segmentId is required" }, { status: 400 });
    }

    await twitchSchedule.deleteSegment(userId, segmentId);
    return NextResponse.json({ ok: true });
  } catch (error) {
    return NextResponse.json(
      { error: error instanceof Error ? error.message : "Unknown error" },
      { status: 400 },
    );
  }
}
