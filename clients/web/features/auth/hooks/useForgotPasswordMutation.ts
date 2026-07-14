import { useMutation } from "@tanstack/react-query";

import { resetPassword } from "@/features/auth/apis";
import { ForgotPasswordForm } from "@/features/auth/schemas";
import { ResetPasswordAPIResponse } from "@/features/auth/types";

// TODO:
// (1) Toast Notification

const useForgotPasswordMutation = () => {
  return useMutation({
    mutationFn: (form: ForgotPasswordForm) =>
      resetPassword({ email: form.email }),
    onSuccess: (data: ResetPasswordAPIResponse | undefined) => {
      // TODO (1)
      //
      // toast.success(data?.email);
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

export default useForgotPasswordMutation;
