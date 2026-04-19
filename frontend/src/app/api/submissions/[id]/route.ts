import { NextRequest, NextResponse } from "next/server";

const BACKEND = "https://azure.abhi.dedyn.io/cfbackend/api";

export async function GET(
  _req: NextRequest,
  { params }: { params: Promise<{ id: string }> },
) {
  try {
    const { id } = await params;

    const backendRes = await fetch(`${BACKEND}/submissions/${id}/`, {
      cache: "no-store",
    });

    const text = await backendRes.text();

    return new NextResponse(text, {
      status: backendRes.status,
      headers: { "Content-Type": "application/json" },
    });
  } catch (err) {
    console.error("[proxy /api/submissions/[id]]", err);
    return new NextResponse("proxy error", { status: 502 });
  }
}
