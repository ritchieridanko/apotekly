"use client";

import React from "react";

import {
  AuthButton,
  OAuthButton,
  RedirectLink,
  TextField,
} from "@/features/auth/components";
import { PASSWORD_MAX_LENGTH } from "@/features/auth/constants";
import { useSignUpForm } from "@/features/auth/hooks";
import { Apple, Google } from "@/shared/assets/images";
import { Logo } from "@/shared/components";
import { ROUTES } from "@/shared/constants";

const SignUp: React.FC = () => {
  const {
    form,
    errors,
    setEmail,
    isEmailAvailable,
    setPassword,
    setConfirmPassword,
    handleSignUp,
    isSigningUp,
  } = useSignUpForm();

  return (
    <div className="w-full flex flex-col justify-center items-center gap-3">
      <div className="w-full flex flex-col justify-center items-center">
        <Logo size="3xl" />
        <p className="w-full font-sans font-normal text-(--dark) text-sm md:text-base text-center">
          Create an account
        </p>
      </div>
      <TextField
        label="Email"
        placeholder="Enter your email here"
        value={form.email}
        onChange={(v) => setEmail(v)}
        isDisabled={isSigningUp}
        isLoading={isSigningUp}
        validLabel={
          isEmailAvailable && !errors.email ? "Email is available" : undefined
        }
        invalidLabel={errors.email}
      />
      <TextField
        label="Password"
        placeholder="Enter your password here"
        value={form.password}
        onChange={(v) => setPassword(v)}
        isDisabled={isSigningUp}
        isLoading={isSigningUp}
        withSecurity
        maxLength={PASSWORD_MAX_LENGTH}
        invalidLabel={errors.password}
      />
      <TextField
        label="Confirm Password"
        placeholder="Enter your password here"
        value={form.confirmPassword}
        onChange={(v) => setConfirmPassword(v)}
        isDisabled={isSigningUp}
        isLoading={isSigningUp}
        withSecurity
        maxLength={PASSWORD_MAX_LENGTH}
        invalidLabel={errors.confirmPassword}
      />
      <div className="w-full font-sans font-normal text-(--dark) text-sm md:text-base text-center">
        <p className="inline">By signing up, I accept the </p>
        <RedirectLink
          label="Terms of Service"
          to={ROUTES.LEGAL.TOS}
          isDisabled={isSigningUp}
          isLoading={isSigningUp}
          customStyle="text-amber-500"
        />
        <p className="inline"> and acknowledge the </p>
        <RedirectLink
          label="Privacy Policy"
          to={ROUTES.LEGAL.PRIVACY_POLICY}
          isDisabled={isSigningUp}
          isLoading={isSigningUp}
          customStyle="text-amber-500"
        />
        {"."}
      </div>
      <AuthButton
        label="Sign Up"
        onClick={() => handleSignUp()}
        isDisabled={isSigningUp}
        loading={{ isLoading: isSigningUp, label: "Signing Up..." }}
      />
      <p className="w-full font-sans font-normal text-gray-400 text-sm md:text-base text-center">
        Or continue with
      </p>
      <OAuthButton
        label="Google"
        icon={<Google className="size-5 md:size-6" />}
        onClick={() => {}}
        isDisabled={isSigningUp}
        isLoading={isSigningUp}
      />
      <OAuthButton
        label="Apple"
        icon={<Apple className="size-5 md:size-6" />}
        onClick={() => {}}
        isDisabled={isSigningUp}
        isLoading={isSigningUp}
      />
      <p className="w-full font-sans font-normal text-(--dark) text-sm md:text-base text-center">
        Already have an account?{" "}
        <RedirectLink
          label="Sign In!"
          to={ROUTES.AUTH.SIGN_IN}
          isDisabled={isSigningUp}
          isLoading={isSigningUp}
        />
      </p>
    </div>
  );
};

export default SignUp;
