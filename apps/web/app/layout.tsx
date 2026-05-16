import type { Metadata } from "next";

import "@/app/globals.css";

export const metadata: Metadata = {
  title: "Cuepoint Control Center",
  description: "Stream planning, bot controls, and overlays for Cuepoint.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body suppressHydrationWarning>{children}</body>
    </html>
  );
}
