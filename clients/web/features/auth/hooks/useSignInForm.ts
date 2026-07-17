import debounce from "lodash/debounce";
import { useCallback, useEffect, useMemo, useState } from "react";
import { ZodSafeParseResult } from "zod";

import { useSignInMutation } from "@/features/auth/hooks";
import { signInSchema, type SignInForm } from "@/features/auth/schemas";
import type { SignInFormErrors } from "@/features/auth/types";
import { setLocalRememberMe } from "@/shared/utils";

const DEBOUNCING_DELAY: number = 300; // 300ms

const useSignInForm = () => {
  const { mutate: signIn, isPending } = useSignInMutation();

  const [form, setForm] = useState<SignInForm>({
    email: "",
    password: "",
  });
  const [errors, setErrors] = useState<SignInFormErrors>({});
  const [rememberMe, setRememberMe] = useState<boolean>(false);

  const validate = useMemo(
    () =>
      debounce(
        async (
          field: keyof SignInForm,
          value: string,
          setErrors: React.Dispatch<React.SetStateAction<SignInFormErrors>>,
        ) => {
          const res: ZodSafeParseResult<string> =
            signInSchema.shape[field].safeParse(value);

          setErrors((prev) => ({
            ...prev,
            [field]: res.success ? undefined : res.error.issues[0]?.message,
          }));
        },
        DEBOUNCING_DELAY,
      ),
    [],
  );

  useEffect(() => {
    return () => validate.cancel();
  }, [validate]);

  const update = useCallback(
    (field: keyof SignInForm, value: string) => {
      setForm((prev) => {
        const nextForm: SignInForm = { ...prev, [field]: value };
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
  const setPassword = useCallback(
    (password: string) => update("password", password),
    [update],
  );

  const handleSignIn = () => {
    const res: ZodSafeParseResult<SignInForm> = signInSchema.safeParse(form);

    if (res.success) {
      setErrors({});
      signIn(
        { form: form, rememberMe: rememberMe },
        {
          onSuccess: () => {
            setLocalRememberMe(rememberMe);
          },
        },
      );
    } else {
      const newErrors: SignInFormErrors = {};

      res.error.issues.forEach((issue) => {
        if (issue.path[0]) {
          newErrors[issue.path[0] as keyof SignInForm] = issue.message;
        }
      });

      setErrors(newErrors);
    }
  };

  return {
    form,
    errors,
    setEmail,
    setPassword,
    rememberMe,
    setRememberMe,
    handleSignIn,
    isSigningIn: isPending,
  };
};

export default useSignInForm;
