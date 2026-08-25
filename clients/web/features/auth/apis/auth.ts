import {
  CheckEmailAvailabilityAPIResponse,
  CheckPasswordResetTokenValidityAPIResponse,
  ConfirmPasswordResetAPIRequest,
  ResetPasswordAPIRequest,
  ResetPasswordAPIResponse,
  SignInAPIRequest,
  SignInAPIResponse,
  SignUpAPIRequest,
  SignUpAPIResponse,
  VerifyEmailAPIRequest,
  VerifyEmailAPIResponse,
} from "@/features/auth/types";
import { api } from "@/shared/libs/api";

const PREFIX: string = "/auth";

const checkEmailAvailability = async (email: string): Promise<boolean> => {
  const res = await api<APIResponse<CheckEmailAvailabilityAPIResponse>>(
    `${PREFIX}/email/available?email=${encodeURIComponent(email)}`,
    { method: "GET" },
  );
  return res.data?.is_available ?? false;
};

const checkPasswordResetTokenValidity = async (
  token: string,
): Promise<boolean> => {
  const res = await api<
    APIResponse<CheckPasswordResetTokenValidityAPIResponse>
  >(`${PREFIX}/password/reset/valid?token=${encodeURIComponent(token)}`, {
    method: "GET",
  });
  return res.data?.is_valid ?? false;
};

const confirmPasswordReset = async (
  form: ConfirmPasswordResetAPIRequest,
): Promise<APIResponse> => {
  return await api<APIResponse>(`${PREFIX}/password/reset/confirm`, {
    method: "POST",
    body: form,
  });
};

const resetPassword = async (
  form: ResetPasswordAPIRequest,
): Promise<APIResponse<ResetPasswordAPIResponse>> => {
  return await api<APIResponse<ResetPasswordAPIResponse>>(
    `${PREFIX}/password/reset`,
    {
      method: "POST",
      body: form,
    },
  );
};

const signIn = async (
  form: SignInAPIRequest,
): Promise<APIResponse<SignInAPIResponse>> => {
  return await api<APIResponse<SignInAPIResponse>>(`${PREFIX}/signin`, {
    method: "POST",
    body: form,
  });
};

const signUp = async (
  form: SignUpAPIRequest,
): Promise<APIResponse<SignUpAPIResponse>> => {
  return await api<APIResponse<SignUpAPIResponse>>(`${PREFIX}/signup`, {
    method: "POST",
    body: form,
  });
};

const verifyEmail = async (
  token: string,
  form: VerifyEmailAPIRequest,
): Promise<APIResponse<VerifyEmailAPIResponse>> => {
  return await api<APIResponse<VerifyEmailAPIResponse>>(
    `${PREFIX}/email/verification/confirm?token=${encodeURIComponent(token)}`,
    {
      method: "POST",
      body: form,
      requiresAuth: true,
    },
  );
};

export {
  checkEmailAvailability,
  checkPasswordResetTokenValidity,
  confirmPasswordReset,
  resetPassword,
  signIn,
  signUp,
  verifyEmail,
};
