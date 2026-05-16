import { NextResponse } from "next/server";

import { requireInternalService } from "@/lib/internal-auth";
import { updateChannelTitle } from "@/lib/twitch";

export const dynamic = "force-dynamic";

export async function POST(
  request: Request,
  { params }: { params: Promise<{ userId: string }> },
) {
  try {
    await requireInternalService();
    const { userId } = await params;
    const body = (await request.json()) as { title?: string };
    if (!body.title?.trim()) {
      return NextResponse.json({ error: "title is required" }, { status: 400 });
    }
    await updateChannelTitle(userId, body.title.trim());
    return NextResponse.json({ ok: true });
  } catch (error) {
    return NextResponse.json(
      { error: error instanceof Error ? error.message : "Unknown error" },
      { status: 400 },
    );
  }
}
