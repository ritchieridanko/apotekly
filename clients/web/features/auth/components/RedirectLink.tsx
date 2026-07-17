"use client";

import { useRouter } from "next/navigation";
import React from "react";

import { Link } from "@/shared/components";

interface RedirectLinkProps {
  label: string;
  to: string;
  isDisabled?: boolean;
  isLoading?: boolean;
  customStyle?: string;
}

const RedirectLink: React.FC<RedirectLinkProps> = ({
  label,
  to,
  isDisabled = false,
  isLoading = false,
  customStyle,
}: RedirectLinkProps) => {
  const router = useRouter();

  return (
    <Link
      onClick={() => router.push(to)}
      isDisabled={isDisabled || isLoading}
      variant={{ color: "primary", isDisabled: isDisabled || isLoading }}
      customStyle={customStyle}
    >
      {label}
    </Link>
  );
};

export default RedirectLink;
