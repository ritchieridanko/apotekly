import debounce from "lodash/debounce";
import { useCallback, useRef, useState } from "react";
import { ZodSafeParseResult } from "zod";

import { useForgotPasswordMutation } from "@/features/auth/hooks";
import {
  forgotPasswordSchema,
  type ForgotPasswordForm,
} from "@/features/auth/schemas";
import type { ForgotPasswordFormErrors } from "@/features/auth/types";

const DEBOUNCING_DELAY: number = 300; // 300ms

const useForgotPasswordForm = () => {
  const { mutate: forgotPassword, isPending } = useForgotPasswordMutation();

  const [form, setForm] = useState<ForgotPasswordForm>({ email: "" });
  const [errors, setErrors] = useState<ForgotPasswordFormErrors>({});

  const validate = useRef(
    debounce(
      async (
        field: keyof ForgotPasswordForm,
        value: string,
        setErrors: React.Dispatch<
          React.SetStateAction<ForgotPasswordFormErrors>
        >,
      ) => {
        const res: ZodSafeParseResult<string> =
          forgotPasswordSchema.shape[field].safeParse(value);

        setErrors((prev) => ({
          ...prev,
          [field]: res.success ? undefined : res.error.issues[0]?.message,
        }));
      },
      DEBOUNCING_DELAY,
    ),
  ).current;

  const update = useCallback(
    (field: keyof ForgotPasswordForm, value: string) => {
      setForm((prev) => {
        const nextForm: ForgotPasswordForm = { ...prev, [field]: value };
        validate(field, value, setErrors);
        return nextForm;
      });
    },
    [validate],
  );
  const setEmail = useCallback(
    (email: string) => update("email", email),
    [update],
  );

  const handleForgotPassword = useCallback(() => {
    const res: ZodSafeParseResult<ForgotPasswordForm> =
      forgotPasswordSchema.safeParse(form);

    if (res.success) {
      setErrors({});
      forgotPassword(form);
    } else {
      const newErrors: ForgotPasswordFormErrors = {};

      res.error.issues.forEach((issue) => {
        if (issue.path[0]) {
          newErrors[issue.path[0] as keyof ForgotPasswordForm] = issue.message;
        }
      });

      setErrors(newErrors);
    }
  }, [form, forgotPassword]);

  return {
    form,
    errors,
    setEmail,
    handleForgotPassword,
    isForgettingPassword: isPending,
  };
};

export default useForgotPasswordForm;
