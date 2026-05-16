import { authHandler } from "@/lib/auth";
import { ensureWebSchema } from "@/lib/db";

export async function GET(request: Request) {
  await ensureWebSchema();
  return authHandler.GET(request);
}

export async function POST(request: Request) {
  await ensureWebSchema();
  return authHandler.POST(request);
}
