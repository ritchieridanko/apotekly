"use client";

import { useRouter } from "next/navigation";
import React from "react";

import { AuthButton, TextField } from "@/features/auth/components";
import { PASSWORD_MAX_LENGTH } from "@/features/auth/constants";
import { useResetPasswordForm } from "@/features/auth/hooks";
import { Logo } from "@/shared/components";
import { ROUTES } from "@/shared/constants";

interface ResetPasswordFormProps {
  token: string;
}

const ResetPasswordForm: React.FC<ResetPasswordFormProps> = ({
  token,
}: ResetPasswordFormProps) => {
  const router = useRouter();

  const {
    form,
    errors,
    setPassword,
    setConfirmPassword,
    handleResetPassword,
    isResettingPassword,
  } = useResetPasswordForm({
    onSuccess: () => {
      router.replace(ROUTES.AUTH.SIGN_IN);
    },
  });

  return (
    <div className="w-full flex flex-col justify-center items-center gap-3">
      <div className="w-full flex flex-col justify-center items-center">
        <Logo size="3xl" />
        <p className="w-full font-sans font-normal text-(--dark) text-sm md:text-base text-center">
          Reset your password
        </p>
      </div>
      <TextField
        label="Password"
        placeholder="Enter your password here"
        value={form.password}
        onChange={(v) => setPassword(v)}
        isDisabled={isResettingPassword}
        isLoading={isResettingPassword}
        withSecurity
        maxLength={PASSWORD_MAX_LENGTH}
        invalidLabel={errors.password}
      />
      <TextField
        label="Confirm Password"
        placeholder="Enter your password here"
        value={form.confirmPassword}
        onChange={(v) => setConfirmPassword(v)}
        isDisabled={isResettingPassword}
        isLoading={isResettingPassword}
        withSecurity
        maxLength={PASSWORD_MAX_LENGTH}
        invalidLabel={errors.confirmPassword}
      />
      <AuthButton
        label="Reset Password"
        onClick={() => handleResetPassword(token)}
        isDisabled={isResettingPassword}
        loading={{
          isLoading: isResettingPassword,
          label: "Resetting Password...",
        }}
      />
    </div>
  );
};

export default ResetPasswordForm;
