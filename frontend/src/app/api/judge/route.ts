import { NextRequest, NextResponse } from "next/server";

const BACKEND = "https://azure.abhi.dedyn.io/cfbackend/api";

export async function POST(req: NextRequest) {
  try {
    const body = await req.json();
    const cookie = req.headers.get("cookie") ?? "";

    const backendRes = await fetch(`${BACKEND}/judge/`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...(cookie ? { Cookie: cookie } : {}),
      },
      body: JSON.stringify(body),
    });

    const text = await backendRes.text();

    return new NextResponse(text, {
      status: backendRes.status,
      headers: { "Content-Type": "application/json" },
    });
  } catch (err) {
    console.error("[proxy /api/judge]", err);
    return new NextResponse("proxy error", { status: 502 });
  }
}
