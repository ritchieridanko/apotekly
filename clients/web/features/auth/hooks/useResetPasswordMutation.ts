import { useMutation } from "@tanstack/react-query";

import { confirmPasswordReset } from "@/features/auth/apis";
import { ResetPasswordForm } from "@/features/auth/schemas";

// TODO:
// (1) Toast Notification

const useResetPasswordMutation = () => {
  return useMutation({
    mutationFn: ({ form, token }: { form: ResetPasswordForm; token: string }) =>
      confirmPasswordReset({ token: token, new_password: form.password }),
    onSuccess: (data: APIResponse) => {
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

export default useResetPasswordMutation;
