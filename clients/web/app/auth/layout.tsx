"use client";

import { useRouter } from "next/navigation";
import React, { useEffect } from "react";

import { useAuthStore } from "@/features/auth/stores";
import { ROUTES } from "@/shared/constants";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }: LayoutProps) => {
  const router = useRouter();
  const { accessToken } = useAuthStore();

  useEffect(() => {
    if (accessToken) {
      router.replace(ROUTES.HOME);
    }
  }, [accessToken, router]);

  if (accessToken) {
    // Guard clause to prevent rendering the auth layout
    // flashing onto the screen while redirecting
    return null;
  }

  return (
    <main
      className="min-h-screen w-full min-w-80 flex justify-center items-center"
      style={{
        backgroundImage: "url(/auth-background-image.webp)",
        backgroundSize: "cover",
        backgroundPosition: "center",
      }}
    >
      <div className="px-5 pt-6 pb-7.5 h-fit w-8/10 md:w-125 min-w-80 bg-(--light) rounded-md ring-8 ring-(--dark)/30">
        {children}
      </div>
    </main>
  );
};

export default Layout;
