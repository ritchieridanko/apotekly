"use client";

import React from "react";

import {
  AuthButton,
  CheckBox,
  OAuthButton,
  RedirectLink,
  TextField,
} from "@/features/auth/components";
import { useSignInForm } from "@/features/auth/hooks";
import { Apple, Google } from "@/shared/assets/images";
import { Logo } from "@/shared/components";
import { ROUTES } from "@/shared/constants";

const SignIn: React.FC = () => {
  const {
    form,
    errors,
    setEmail,
    setPassword,
    rememberMe,
    setRememberMe,
    handleSignIn,
    isSigningIn,
  } = useSignInForm();

  return (
    <div className="w-full flex flex-col justify-center items-center gap-3">
      <div className="w-full flex flex-col justify-center items-center">
        <Logo size="3xl" />
        <p className="w-full font-sans font-normal text-(--dark) text-sm md:text-base text-center">
          Sign in to continue
        </p>
      </div>
      <TextField
        label="Email"
        placeholder="Enter your email here"
        value={form.email}
        onChange={(v) => setEmail(v)}
        isDisabled={isSigningIn}
        isLoading={isSigningIn}
        invalidLabel={errors.email}
      />
      <TextField
        label="Password"
        placeholder="Enter your password here"
        value={form.password}
        onChange={(v) => setPassword(v)}
        isDisabled={isSigningIn}
        isLoading={isSigningIn}
        withSecurity
        invalidLabel={errors.password}
      />
      <div className="w-full flex flex-row justify-between items-center font-sans font-normal text-(--dark) text-sm md:text-base text-center">
        <CheckBox
          label="Remember Me"
          isChecked={rememberMe}
          onChange={(v) => setRememberMe(v)}
          isDisabled={isSigningIn}
          isLoading={isSigningIn}
        />
        <RedirectLink
          label="Forgot Password?"
          to={ROUTES.AUTH.FORGOT_PASSWORD}
          isDisabled={isSigningIn}
          isLoading={isSigningIn}
        />
      </div>
      <AuthButton
        label="Sign In"
        onClick={() => handleSignIn()}
        isDisabled={isSigningIn}
        loading={{ isLoading: isSigningIn, label: "Signing In..." }}
      />
      <p className="w-full font-sans font-normal text-gray-400 text-sm md:text-base text-center">
        Or continue with
      </p>
      <OAuthButton
        label="Google"
        icon={<Google className="size-5 md:size-6" />}
        onClick={() => {}}
        isDisabled={isSigningIn}
        isLoading={isSigningIn}
      />
      <OAuthButton
        label="Apple"
        icon={<Apple className="size-5 md:size-6" />}
        onClick={() => {}}
        isDisabled={isSigningIn}
        isLoading={isSigningIn}
      />
      <p className="w-full font-sans font-normal text-(--dark) text-sm md:text-base text-center">
        Don't have an account?{" "}
        <RedirectLink
          label="Sign Up!"
          to={ROUTES.AUTH.SIGN_UP}
          isDisabled={isSigningIn}
          isLoading={isSigningIn}
        />
      </p>
    </div>
  );
};

export default SignIn;
