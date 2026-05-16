import { NextResponse } from "next/server";

import { requireInternalService } from "@/lib/internal-auth";
import { getStreamState } from "@/lib/twitch";

export const dynamic = "force-dynamic";

export async function GET(
  _request: Request,
  { params }: { params: Promise<{ userId: string }> },
) {
  try {
    await requireInternalService();
    const { userId } = await params;
    return NextResponse.json(await getStreamState(userId));
  } catch (error) {
    return NextResponse.json(
      { error: error instanceof Error ? error.message : "Unknown error" },
      { status: 400 },
    );
  }
}
