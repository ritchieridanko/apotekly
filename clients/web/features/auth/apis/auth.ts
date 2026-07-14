import {
  CheckEmailAvailabilityAPIResponse,
  SignInAPIRequest,
  SignInAPIResponse,
  SignUpAPIRequest,
  SignUpAPIResponse,
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

const signIn = async (
  form: SignInAPIRequest,
): Promise<SignInAPIResponse | undefined> => {
  const res = await api<APIResponse<SignInAPIResponse>>(`${PREFIX}/signin`, {
    method: "POST",
    body: form,
  });
  return res.data;
};

const signUp = async (
  form: SignUpAPIRequest,
): Promise<SignUpAPIResponse | undefined> => {
  const res = await api<APIResponse<SignUpAPIResponse>>(`${PREFIX}/signup`, {
    method: "POST",
    body: form,
  });
  return res.data;
};

export { checkEmailAvailability, signIn, signUp };
