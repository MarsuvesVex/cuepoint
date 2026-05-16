import { NextRequest, NextResponse } from "next/server";

import { requireInternalService } from "@/lib/internal-auth";
import { getOperatorProfile, twitchVideos } from "@/lib/twitch";

export const dynamic = "force-dynamic";

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ userId: string }> },
) {
  try {
    await requireInternalService();
    const { userId } = await params;
    const ids = request.nextUrl.searchParams.getAll("id");
    const first = request.nextUrl.searchParams.get("first");
    const after = request.nextUrl.searchParams.get("after") ?? undefined;
    const profile = await getOperatorProfile(userId);

    return NextResponse.json(
      await twitchVideos.list({
        userId,
        ids: ids.length > 0 ? ids : undefined,
        channelUserId: ids.length === 0 ? profile?.twitch_broadcaster_id : undefined,
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
