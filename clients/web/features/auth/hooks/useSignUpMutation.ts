import { useMutation } from "@tanstack/react-query";
import Cookies from "js-cookie";

import { signUp } from "@/features/auth/apis";
import { SignUpForm } from "@/features/auth/schemas";
import { useAuthStore } from "@/features/auth/stores";
import { SignUpAPIResponse } from "@/features/auth/types";

// TODO:
// (1) Toast Notification

const ENV: string = process.env.NEXT_PUBLIC_APP_ENV ?? "dev";

const useSignUpMutation = () => {
  const { setAuth, setAccessToken } = useAuthStore();

  return useMutation({
    mutationFn: (form: SignUpForm) =>
      signUp({ email: form.email, password: form.password }),
    onSuccess: (data: SignUpAPIResponse | undefined) => {
      if (data?.auth) {
        setAuth({
          email: data.auth.email,
          role: data.auth.role,
          isEmailVerified: data.auth.is_email_verified,
        });
      }
      if (data?.access_token?.token) {
        setAccessToken(data.access_token.token);

        const seconds: number = data.access_token.expires_in_seconds;
        const expiryDate: Date = new Date(Date.now() + seconds * 1000);

        Cookies.set("access_token", data.access_token.token, {
          expires: expiryDate,
          secure: ENV === "prod",
          sameSite: "Strict",
        });
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

export default useSignUpMutation;
