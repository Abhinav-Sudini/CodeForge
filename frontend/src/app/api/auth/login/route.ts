import { NextRequest, NextResponse } from "next/server";

const BACKEND = "https://azure.abhi.dedyn.io/cfbackend/auth";

export async function POST(req: NextRequest) {
  try {
    const body = await req.json();

    const backendRes = await fetch(`${BACKEND}/login/`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });

    const text = await backendRes.text();

    const res = new NextResponse(text, {
      status: backendRes.status,
      headers: { "Content-Type": "text/plain" },
    });

    // Forward cookies safely to localhost
    const setCookies = backendRes.headers.getSetCookie();
    for (let cookie of setCookies) {
      cookie = cookie.replace(/Domain=[^;]+;?\s*/i, "");
      cookie = cookie.replace(/Secure;?\s*/i, "");
      cookie = cookie.replace(/SameSite=[^;]+;?\s*/i, "");
      cookie = cookie.replace(/Path=[^;]+;?\s*/i, "");
      cookie = cookie.trim();
      if (!cookie.endsWith(';')) cookie += ';';
      cookie += " Path=/; SameSite=Lax";
      res.headers.append("Set-Cookie", cookie);
    }

    return res;
  } catch (err) {
    console.error("[proxy /api/auth/login]", err);
    return new NextResponse("proxy error", { status: 502 });
  }
}
