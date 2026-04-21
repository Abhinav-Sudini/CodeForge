import { NextRequest, NextResponse } from "next/server";

const BACKEND = "https://azure.abhi.dedyn.io/cfbackend/api";

export async function GET(
  req: NextRequest,
  { params }: { params: Promise<{ qid: string }> },
) {
  try {
    const { qid } = await params;
    const cookie = req.headers.get("cookie") ?? "";

    const backendRes = await fetch(
      `${BACKEND}/question/submissions/${qid}/`,
      {
        cache: "no-store",
        headers: {
          ...(cookie ? { Cookie: cookie } : {}),
        },
      },
    );

    const text = await backendRes.text();

    return new NextResponse(text, {
      status: backendRes.status,
      headers: { "Content-Type": "application/json" },
    });
  } catch (err) {
    console.error("[proxy /api/question-submissions/[qid]]", err);
    return new NextResponse("proxy error", { status: 502 });
  }
}
