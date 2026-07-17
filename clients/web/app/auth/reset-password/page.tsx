import { QueryClient } from "@tanstack/react-query";
import React from "react";

import ResetPasswordForm from "./form";
import { checkPasswordResetTokenValidity } from "@/features/auth/apis";
import { AuthFeedbackTile } from "@/features/auth/components";
import { XCircle } from "@/shared/assets/icons";

const TOKEN_VALIDITY_CACHE_TIMEOUT: number = 30 * 1000; // 30s

interface ResetPasswordPageProps {
  searchParams: Promise<{ token?: string }>;
}

const ResetPasswordPage: React.FC<ResetPasswordPageProps> = async ({
  searchParams,
}: ResetPasswordPageProps) => {
  const token: string | undefined = (await searchParams).token;

  if (!token) {
    return (
      <AuthFeedbackTile
        message="Missing Token"
        icon={XCircle}
        variant="error"
      />
    );
  }

  const queryClient: QueryClient = new QueryClient();

  try {
    const isValid: boolean = await queryClient.fetchQuery({
      queryKey: ["password-reset-token-validity", token],
      queryFn: () => checkPasswordResetTokenValidity(token),
      staleTime: TOKEN_VALIDITY_CACHE_TIMEOUT,
    });

    if (!isValid) {
      return (
        <AuthFeedbackTile
          message="Invalid Token"
          icon={XCircle}
          variant="error"
        />
      );
    }
  } catch {
    return (
      <AuthFeedbackTile
        message="Internal Server Error. Try Again Later!"
        icon={XCircle}
        variant="error"
      />
    );
  }

  return <ResetPasswordForm token={token} />;
};

export default ResetPasswordPage;
