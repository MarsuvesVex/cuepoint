import { NextRequest, NextResponse } from "next/server";

import { requireInternalService } from "@/lib/internal-auth";
import { twitchMarkers } from "@/lib/twitch";

export const dynamic = "force-dynamic";

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ userId: string }> },
) {
  try {
    await requireInternalService();
    const { userId } = await params;
    const videoId = request.nextUrl.searchParams.get("videoId") ?? undefined;
    const first = request.nextUrl.searchParams.get("first");
    const after = request.nextUrl.searchParams.get("after") ?? undefined;
    const before = request.nextUrl.searchParams.get("before") ?? undefined;

    return NextResponse.json(
      await twitchMarkers.list({
        userId,
        videoId,
        first: first ? Number(first) : undefined,
        after,
        before,
      }),
    );
  } catch (error) {
    return NextResponse.json(
      { error: error instanceof Error ? error.message : "Unknown error" },
      { status: 400 },
    );
  }
}

export async function POST(
  request: Request,
  { params }: { params: Promise<{ userId: string }> },
) {
  try {
    await requireInternalService();
    const { userId } = await params;
    const body = (await request.json()) as { description?: string };
    return NextResponse.json(await twitchMarkers.create(userId, body.description));
  } catch (error) {
    return NextResponse.json(
      { error: error instanceof Error ? error.message : "Unknown error" },
      { status: 400 },
    );
  }
}
