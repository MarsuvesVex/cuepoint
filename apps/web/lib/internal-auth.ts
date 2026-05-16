import { headers } from "next/headers";

export async function requireInternalService() {
  const headerName =
    process.env.CUEPOINT_INTERNAL_HEADER_NAME ?? "X-Cuepoint-Internal-Token";
  const token = process.env.CUEPOINT_INTERNAL_SERVICE_TOKEN ?? "";
  const incoming = (await headers()).get(headerName);

  if (token && incoming !== token) {
    throw new Error("UNAUTHORIZED");
  }
}
