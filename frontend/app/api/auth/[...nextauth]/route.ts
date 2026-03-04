import { NextResponse } from "next/server";

function disabledResponse() {
  return NextResponse.json(
    {
      error: "NextAuth is disabled. Use Go backend auth at /auth/login.",
    },
    { status: 410 }
  );
}

export async function GET() {
  return disabledResponse();
}

export async function POST() {
  return disabledResponse();
}
