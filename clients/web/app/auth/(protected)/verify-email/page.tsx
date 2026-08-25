import React from "react";

import VerifyEmailTrigger from "./trigger";
import { AuthFeedbackTile } from "@/features/auth/components";
import { XCircle } from "@/shared/assets/icons";

interface VerifyEmailPageProps {
  searchParams: Promise<{ token?: string }>;
}

const VerifyEmailPage: React.FC<VerifyEmailPageProps> = async ({
  searchParams,
}: VerifyEmailPageProps) => {
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

  return <VerifyEmailTrigger token={token} />;
};

export default VerifyEmailPage;
