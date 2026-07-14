"use client";

import React from "react";

import {
  AuthButton,
  RedirectLink,
  TextField,
} from "@/features/auth/components";
import { useForgotPasswordForm } from "@/features/auth/hooks";
import { Logo } from "@/shared/components";
import { ROUTES } from "@/shared/constants";

const ForgotPassword: React.FC = () => {
  const { form, errors, setEmail, handleForgotPassword, isForgettingPassword } =
    useForgotPasswordForm();

  return (
    <div className="w-full flex flex-col justify-center items-center gap-3">
      <div className="w-full flex flex-col justify-center items-center">
        <Logo size="3xl" />
        <p className="w-full font-sans font-normal text-(--dark) text-sm md:text-base text-center">
          Forgot your password?
        </p>
      </div>
      <TextField
        label="Email"
        placeholder="Enter your email here"
        value={form.email}
        onChange={(v) => setEmail(v)}
        isDisabled={isForgettingPassword}
        isLoading={isForgettingPassword}
        invalidLabel={errors.email}
      />
      <AuthButton
        label="Reset Password"
        onClick={() => handleForgotPassword()}
        isDisabled={isForgettingPassword}
        loading={{
          isLoading: isForgettingPassword,
          label: "Resetting Password...",
        }}
      />
      <p className="w-full font-sans font-normal text-(--dark) text-sm md:text-base text-center">
        Don't have an account?{" "}
        <RedirectLink
          label="Sign Up!"
          to={ROUTES.AUTH.SIGN_UP}
          isDisabled={isForgettingPassword}
          isLoading={isForgettingPassword}
        />
      </p>
    </div>
  );
};

export default ForgotPassword;
