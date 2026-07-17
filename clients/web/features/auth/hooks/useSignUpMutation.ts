import { useMutation } from "@tanstack/react-query";

import { signUp } from "@/features/auth/apis";
import { SignUpForm } from "@/features/auth/schemas";
import { useAuthStore } from "@/features/auth/stores";
import { SignUpAPIResponse } from "@/features/auth/types";
import { APP_ENV, setCookieAccessToken } from "@/shared/utils";

// TODO:
// (1) Toast Notification

const useSignUpMutation = () => {
  const { setAuth, setAccessToken } = useAuthStore();

  return useMutation({
    mutationFn: (form: SignUpForm) =>
      signUp({ email: form.email, password: form.password }),
    onSuccess: (data: APIResponse<SignUpAPIResponse>) => {
      if (data.data?.auth) {
        setAuth({
          email: data.data.auth.email,
          role: data.data.auth.role,
          isEmailVerified: data.data.auth.is_email_verified,
        });
      }
      if (data.data?.access_token) {
        const token: string = data.data.access_token.token;
        setAccessToken(token);

        const seconds: number = data.data.access_token.expires_in_seconds;
        const expiryDate: Date = new Date(Date.now() + seconds * 1000);
        const attributes: Cookies.CookieAttributes = {
          expires: expiryDate,
          secure: APP_ENV === "prod",
          sameSite: "Strict",
        };

        setCookieAccessToken(token, attributes);
      }

      // TODO (1)
      //
      // toast.success(data.message);
    },
    onError: (error: Error) => {
      // TODO (1)
      //
      // if (error instanceof APIError) {
      //   toast.error(error.message);
      // } else {
      //   toast.error("Internal server error");
      // }
    },
  });
};

export default useSignUpMutation;
