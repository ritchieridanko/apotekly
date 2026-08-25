import debounce from "lodash/debounce";
import { useCallback, useEffect, useMemo, useState } from "react";
import { ZodSafeParseResult } from "zod";

import { useResetPasswordMutation } from "@/features/auth/hooks";
import {
  resetPasswordSchema,
  type ResetPasswordForm,
} from "@/features/auth/schemas";
import type { ResetPasswordFormErrors } from "@/features/auth/types";

const DEBOUNCING_DELAY: number = 300; // 300ms

interface UseResetPasswordFormProps {
  onSuccess?: () => void;
}

const useResetPasswordForm = ({
  onSuccess,
}: UseResetPasswordFormProps = {}) => {
  const { mutate: resetPassword, isPending } = useResetPasswordMutation();

  const [form, setForm] = useState<ResetPasswordForm>({
    password: "",
    confirmPassword: "",
  });
  const [errors, setErrors] = useState<ResetPasswordFormErrors>({});

  const validate = useMemo(
    () =>
      debounce(
        async (
          field: keyof ResetPasswordForm,
          value: string,
          currentPasswordValue: string,
          setErrors: React.Dispatch<
            React.SetStateAction<ResetPasswordFormErrors>
          >,
        ) => {
          const res: ZodSafeParseResult<string> =
            resetPasswordSchema.shape[field].safeParse(value);

          if (res.success && field === "confirmPassword") {
            setErrors((prev) => ({
              ...prev,
              confirmPassword:
                value === currentPasswordValue
                  ? undefined
                  : "Passwords do not match",
            }));
          } else {
            setErrors((prev) => ({
              ...prev,
              [field]: res.success ? undefined : res.error.issues[0]?.message,
            }));
          }
        },
        DEBOUNCING_DELAY,
      ),
    [],
  );

  useEffect(() => {
    return () => validate.cancel();
  }, [validate]);

  const update = useCallback(
    (field: keyof ResetPasswordForm, value: string) => {
      setForm((prev) => {
        const nextForm: ResetPasswordForm = { ...prev, [field]: value };
        validate(field, value, nextForm.password, setErrors);
        return nextForm;
      });
    },
    [validate],
  );
  const setPassword = useCallback(
    (password: string) => update("password", password),
    [update],
  );
  const setConfirmPassword = useCallback(
    (password: string) => update("confirmPassword", password),
    [update],
  );

  const handleResetPassword = (token: string) => {
    const res: ZodSafeParseResult<ResetPasswordForm> =
      resetPasswordSchema.safeParse(form);

    if (res.success && form.password !== form.confirmPassword) {
      setErrors((prev) => ({
        ...prev,
        confirmPassword: "Passwords do not match",
      }));
      return;
    } else if (res.success) {
      setErrors({});
      resetPassword(
        { form: form, token: token },
        {
          onSuccess: () => {
            onSuccess?.();
          },
        },
      );
    } else {
      const newErrors: ResetPasswordFormErrors = {};

      res.error.issues.forEach((issue) => {
        if (issue.path[0]) {
          newErrors[issue.path[0] as keyof ResetPasswordForm] = issue.message;
        }
      });

      setErrors(newErrors);
    }
  };

  return {
    form,
    errors,
    setPassword,
    setConfirmPassword,
    handleResetPassword,
    isResettingPassword: isPending,
  };
};

export default useResetPasswordForm;
