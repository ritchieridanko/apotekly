"use client";

import { useRouter } from "next/navigation";
import React, { useEffect, useRef, useState } from "react";

import { AuthFeedbackTile } from "@/features/auth/components";
import { useVerifyEmailMutation } from "@/features/auth/hooks";
import { useAuthStore } from "@/features/auth/stores";
import { Spinner, XCircle } from "@/shared/assets/icons";
import { ROUTES } from "@/shared/constants";
import { APIError } from "@/shared/libs";
import { getLocalRememberMe } from "@/shared/utils";

const UNAUTHENTICATED_ERRORS: string[] = [
  "Invalid session",
  "Session expired",
  "Unauthenticated",
];

interface VerifyEmailTriggerProps {
  token: string;
}

const VerifyEmailTrigger: React.FC<VerifyEmailTriggerProps> = ({
  token,
}: VerifyEmailTriggerProps) => {
  const router = useRouter();
  const accessToken: string | null = useAuthStore((s) => s.accessToken);

  const { mutate: verifyEmail, isPending: isVerifyingEmail } =
    useVerifyEmailMutation();

  // Prevent strict mode double-triggering in development environments
  const hasTriggered = useRef<boolean>(false);
  const shouldRetryAfterAuth = useRef<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (hasTriggered.current && !shouldRetryAfterAuth.current) return;

    hasTriggered.current = true;
    shouldRetryAfterAuth.current = false;

    verifyEmail(
      { token: token, rememberMe: getLocalRememberMe() },
      {
        onSuccess: () => {
          router.replace(ROUTES.HOME);
        },
        onError: (error) => {
          if (error instanceof APIError) {
            if (
              error.status === 401 &&
              UNAUTHENTICATED_ERRORS.includes(error.message)
            ) {
              const callback: string = encodeURIComponent(
                `${ROUTES.AUTH.VERIFY_EMAIL}?token=${token}`,
              );

              shouldRetryAfterAuth.current = true;
              router.replace(`${ROUTES.AUTH.SIGN_IN}?callback=${callback}`);
            } else {
              setError(
                error.status === 400
                  ? "Invalid Token"
                  : "Internal Server Error. Try Again Later!",
              );
            }
          } else {
            setError("Something Went Wrong. Try Again Later!");
          }
        },
      },
    );
  }, [router, accessToken, token]);

  if (error) {
    return <AuthFeedbackTile message={error} icon={XCircle} variant="error" />;
  }
  if (isVerifyingEmail) {
    return (
      <div className="mt-2 w-full flex flex-col justify-center items-center gap-2 text-(--primary)">
        <Spinner size={12} strokeWidth={3.5} />
        <p className="font-sans font-medium text-sm md:text-base text-center select-none">
          Verifying Your Email...
        </p>
      </div>
    );
  }

  return null;
};

export default VerifyEmailTrigger;
