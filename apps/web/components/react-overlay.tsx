"use client";

import { useEffect, useState } from "react";

type OverlayData = {
  active: boolean;
  segmentTitle: string;
  videoTitle: string;
  creatorName: string;
  url: string;
};

export function ReactOverlay({
  overlayPublicID,
  mode,
}: {
  overlayPublicID: string;
  mode: "persistent" | "periodic";
}) {
  const [data, setData] = useState<OverlayData | null>(null);
  const [visible, setVisible] = useState(mode === "persistent");

  useEffect(() => {
    let cancelled = false;

    const load = async () => {
      const response = await fetch(
        `/api/public/overlay/react?user=${encodeURIComponent(overlayPublicID)}`,
        { cache: "no-store" },
      );
      if (!response.ok) return;
      const json = (await response.json()) as OverlayData;
      if (!cancelled) {
        setData(json);
      }
    };

    void load();
    const poll = window.setInterval(load, 15000);

    let periodicTimer = 0;
    if (mode === "periodic") {
      periodicTimer = window.setInterval(() => {
        setVisible(true);
        window.setTimeout(() => setVisible(false), 20_000);
      }, 10 * 60 * 1000);
      setVisible(true);
      window.setTimeout(() => setVisible(false), 20_000);
    }

    return () => {
      cancelled = true;
      window.clearInterval(poll);
      if (periodicTimer) window.clearInterval(periodicTimer);
    };
  }, [mode, overlayPublicID]);

  const shouldRender =
    Boolean(data?.active) && (mode === "persistent" ? true : visible);

  if (!shouldRender || !data) {
    return <div className="hidden" />;
  }

  return (
    <div className="flex min-h-screen items-end p-8">
      <div className="max-w-xl rounded-[1.75rem] border border-white/20 bg-black/75 px-6 py-5 text-white shadow-2xl backdrop-blur">
        <p className="text-xs uppercase tracking-[0.22em] text-white/60">
          Now reacting
        </p>
        <p className="mt-2 text-2xl font-semibold">{data.videoTitle || data.segmentTitle}</p>
        <p className="mt-1 text-base text-white/75">{data.creatorName}</p>
        <p className="mt-3 text-sm text-white/60">{data.url}</p>
      </div>
    </div>
  );
}
