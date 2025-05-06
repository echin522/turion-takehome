import { proxyTelemetry } from "@/app/api/telemetry/telemetryProxy"
import { NextRequest } from "next/server"

export async function GET(req: NextRequest) {
  return proxyTelemetry(req, {
    upstreamPath: "/api/v1/telemetry/anomaly",
  })
}
