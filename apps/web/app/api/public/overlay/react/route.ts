import { NextRequest, NextResponse } from "next/server";

import { getOverlayData } from "@/lib/data";

export async function GET(request: NextRequest) {
  const overlayPublicID = request.nextUrl.searchParams.get("user");
  if (!overlayPublicID) {
    return NextResponse.json({ error: "user is required" }, { status: 400 });
  }

  const data = await getOverlayData(overlayPublicID);
  return NextResponse.json({
    active: Boolean(data && data.segment_type === "react"),
    segmentTitle: data?.segment_title ?? "",
    videoTitle: data?.youtube_video_title ?? "",
    creatorName: data?.youtube_creator_name ?? "",
    url: data?.youtube_canonical_url ?? "",
  });
}
