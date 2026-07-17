import { useMutation } from "@tanstack/react-query";

import { signIn } from "@/features/auth/apis";
import { SignInForm } from "@/features/auth/schemas";
import { useAuthStore } from "@/features/auth/stores";
import { SignInAPIResponse } from "@/features/auth/types";
import { APP_ENV, setCookieAccessToken } from "@/shared/utils";

// TODO:
// (1) Toast Notification

const useSignInMutation = () => {
  const { setAuth, setAccessToken } = useAuthStore();

  return useMutation({
    mutationFn: ({
      form,
      rememberMe,
    }: {
      form: SignInForm;
      rememberMe: boolean;
    }) =>
      signIn({
        email: form.email,
        password: form.password,
        remember_me: rememberMe,
      }),
    onSuccess: (data: APIResponse<SignInAPIResponse>, { rememberMe }) => {
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

        const attributes: Cookies.CookieAttributes = {
          secure: APP_ENV === "prod",
          sameSite: "Strict",
        };
        if (rememberMe) {
          const seconds: number = data.data.access_token.expires_in_seconds;
          attributes.expires = new Date(Date.now() + seconds * 1000);
        }

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

export default useSignInMutation;
