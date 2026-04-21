import { NextRequest, NextResponse } from "next/server";

const BACKEND = "https://azure.abhi.dedyn.io/cfbackend/auth";

export async function POST(req: NextRequest) {
  try {
    const body = await req.json();

    const backendRes = await fetch(`${BACKEND}/register/`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });

    const text = await backendRes.text();

    return new NextResponse(text, {
      status: backendRes.status,
      headers: { "Content-Type": "text/plain" },
    });
  } catch (err) {
    console.error("[proxy /api/auth/register]", err);
    return new NextResponse("proxy error", { status: 502 });
  }
}
