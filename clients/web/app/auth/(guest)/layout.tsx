"use client";

import { useRouter, useSearchParams } from "next/navigation";
import React, { useEffect } from "react";

import { useAuthStore } from "@/features/auth/stores";
import { ROUTES } from "@/shared/constants";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }: LayoutProps) => {
  const router = useRouter();
  const hasCallback: boolean = !!useSearchParams().get("callback");
  const accessToken: string | null = useAuthStore((s) => s.accessToken);

  useEffect(() => {
    if (accessToken && !hasCallback) {
      router.replace(ROUTES.HOME);
    }
  }, [accessToken, hasCallback, router]);

  if (accessToken && !hasCallback) {
    // Guard clause to prevent rendering the auth layout
    // flashing onto the screen while redirecting
    return null;
  }

  return <>{children}</>;
};

export default Layout;
