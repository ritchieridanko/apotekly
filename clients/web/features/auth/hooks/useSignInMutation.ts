import { useMutation } from "@tanstack/react-query";
import Cookies from "js-cookie";

import { signIn } from "@/features/auth/apis";
import { SignInForm } from "@/features/auth/schemas";
import { useAuthStore } from "@/features/auth/stores";
import { SignInAPIResponse } from "@/features/auth/types";

// TODO:
// (1) Toast Notification

const ENV: string = process.env.NEXT_PUBLIC_APP_ENV ?? "dev";

const useSignInMutation = () => {
  const { setAuth, setAccessToken } = useAuthStore();

  return useMutation({
    mutationFn: ({
      form,
      persisted,
    }: {
      form: SignInForm;
      persisted: boolean;
    }) => signIn({ email: form.email, password: form.password }),
    onSuccess: (data: SignInAPIResponse | undefined, { persisted }) => {
      if (data?.auth) {
        setAuth({
          email: data.auth.email,
          role: data.auth.role,
          isEmailVerified: data.auth.is_email_verified,
        });
      }
      if (data?.access_token?.token) {
        setAccessToken(data.access_token.token);

        const options: Cookies.CookieAttributes = {
          secure: ENV === "prod",
          sameSite: "strict",
        };

        if (persisted) {
          const seconds: number = data.access_token.expires_in_seconds;
          options.expires = new Date(Date.now() + seconds * 1000);
        }

        Cookies.set("access_token", data.access_token.token, options);
      }

      // TODO (1)
      //
      // toast.success("...");
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
