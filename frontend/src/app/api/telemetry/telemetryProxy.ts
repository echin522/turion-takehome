// lib/telemetryProxy.ts
import { NextRequest, NextResponse } from "next/server"

interface ProxyOpts {
  upstreamPath: string        // e.g. "/api/v1/telemetry"
  requireRange?: boolean      
}

export async function proxyTelemetry(
  req: NextRequest,
  opts: ProxyOpts
): Promise<NextResponse> {
  const base = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8090"
  if (!base) {
    return NextResponse.json(
      { error: "`NEXT_PUBLIC_API_URL` is not configured" },
      { status: 500 }
    )
  }

  const url = new URL(opts.upstreamPath, base)

  if (opts.requireRange ?? true) {
    const start = req.nextUrl.searchParams.get("start_time")
    const end   = req.nextUrl.searchParams.get("end_time")
    if (!start || !end) {
      return NextResponse.json(
        { error: "`start_time` and `end_time` are required" },
        { status: 400 }
      )
    }
    url.searchParams.set("start_time", start)
    url.searchParams.set("end_time", end)
  }


  const resp = await fetch(url)
  const data = await resp.text()
  return new NextResponse(data, {
    status: resp.status,
    headers: { "content-type": "application/json" },
  })
}
