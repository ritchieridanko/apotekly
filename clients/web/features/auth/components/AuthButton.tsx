"use client";

import React from "react";

import { Spinner } from "@/shared/assets/icons";
import { Button } from "@/shared/components";

interface AuthButtonProps {
  label: string;
  onClick: () => void;
  isDisabled?: boolean;
  loading: {
    isLoading?: boolean;
    label?: string;
  };
}

const AuthButton: React.FC<AuthButtonProps> = ({
  label,
  onClick,
  isDisabled = false,
  loading = {
    isLoading: false,
    label,
  },
}: AuthButtonProps) => {
  return (
    <Button
      onClick={onClick}
      isDisabled={isDisabled || loading.isLoading}
      variant={{ color: "primary", size: "sm" }}
      customStyle="focus:ring-4 focus:ring-(--primary)/30"
    >
      {loading.isLoading ? (
        <div className="w-full flex flex-row justify-center items-center gap-2 md:gap-2.5">
          <Spinner size={6} strokeWidth={3.5} />
          <p className="text-base md:text-lg">{loading.label}</p>
        </div>
      ) : (
        <p className="font-medium text-base md:text-lg">{label}</p>
      )}
    </Button>
  );
};

export default AuthButton;
