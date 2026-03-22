import { NextResponse } from "next/server";

export async function GET() {
  // In production: validate session cookie, return user info
  // For now return unauthenticated response
  return NextResponse.json(null, { status: 401 });
}
