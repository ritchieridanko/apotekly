import { NextResponse, type NextRequest } from "next/server";

export function middleware(request: NextRequest) {
  const token: string | undefined = request.cookies.get("access_token")?.value;
  const { pathname } = request.nextUrl;

  const isAuthPage: boolean = pathname.startsWith("/auth");
  const isVerifyEmailPage: boolean = pathname === "/auth/verify-email";

  if (isAuthPage && token && !isVerifyEmailPage) {
    return NextResponse.redirect(new URL("/", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/auth/:path*"],
};
