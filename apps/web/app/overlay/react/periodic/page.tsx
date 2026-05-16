import { ReactOverlay } from "@/components/react-overlay";

export const dynamic = "force-dynamic";

export default async function PeriodicReactOverlayPage({
  searchParams,
}: {
  searchParams: Promise<{ user?: string }>;
}) {
  const params = await searchParams;

  return <ReactOverlay overlayPublicID={params.user ?? ""} mode="periodic" />;
}
