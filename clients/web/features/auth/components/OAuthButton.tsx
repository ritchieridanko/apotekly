"use client";

import React from "react";

import { Button } from "@/shared/components";

interface OAuthButtonProps {
  label: string;
  icon: React.ReactNode;
  onClick: () => void;
  isDisabled?: boolean;
  isLoading?: boolean;
}

const OAuthButton: React.FC<OAuthButtonProps> = ({
  label,
  icon,
  onClick,
  isDisabled = false,
  isLoading = false,
}: OAuthButtonProps) => {
  return (
    <Button
      onClick={onClick}
      isDisabled={isDisabled || isLoading}
      variant={{ color: "secondary", size: "sm" }}
      customStyle="focus:ring-4 focus:ring-(--primary)/30"
    >
      <div className="w-full flex flex-row justify-center items-center gap-1.5 md:gap-2">
        {icon}
        <p className={`text-base md:text-lg ${!isLoading && "font-medium"}`}>
          {label}
        </p>
      </div>
    </Button>
  );
};

export default OAuthButton;
