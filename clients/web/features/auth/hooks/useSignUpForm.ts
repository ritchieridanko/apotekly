import { QueryClient, useQueryClient } from "@tanstack/react-query";
import debounce from "lodash/debounce";
import { useCallback, useEffect, useMemo, useState } from "react";
import { ZodSafeParseResult } from "zod";

import { checkEmailAvailability } from "@/features/auth/apis";
import { useSignUpMutation } from "@/features/auth/hooks";
import { signUpSchema, type SignUpForm } from "@/features/auth/schemas";
import type { SignUpFormErrors } from "@/features/auth/types";
import { setLocalRememberMe } from "@/shared/utils";

// TODO:
// (1) Toast Notification

const DEBOUNCING_DELAY: number = 300; // 300ms
const EMAIL_AVAILABILITY_CACHE_TIMEOUT: number = 30 * 1000; // 30s

const useSignUpForm = () => {
  const queryClient: QueryClient = useQueryClient();
  const { mutate: signUp, isPending } = useSignUpMutation();

  const [form, setForm] = useState<SignUpForm>({
    email: "",
    password: "",
    confirmPassword: "",
  });
  const [errors, setErrors] = useState<SignUpFormErrors>({});
  const [isEmailAvailable, setIsEmailAvailable] = useState<boolean>(false);

  const validate = useMemo(
    () =>
      debounce(
        async (
          field: keyof SignUpForm,
          value: string,
          currentPasswordValue: string,
          queryClient: QueryClient,
          setErrors: React.Dispatch<React.SetStateAction<SignUpFormErrors>>,
          setIsEmailAvailable: React.Dispatch<React.SetStateAction<boolean>>,
        ) => {
          const res: ZodSafeParseResult<string> =
            signUpSchema.shape[field].safeParse(value);

          if (res.success && field === "confirmPassword") {
            setErrors((prev) => ({
              ...prev,
              confirmPassword:
                value === currentPasswordValue
                  ? undefined
                  : "Passwords do not match",
            }));
          } else if (res.success && field === "email") {
            setErrors((prev) => ({ ...prev, email: undefined }));

            try {
              const isAvailable: boolean = await queryClient.fetchQuery({
                queryKey: ["email-availability", value],
                queryFn: () => checkEmailAvailability(value),
                staleTime: EMAIL_AVAILABILITY_CACHE_TIMEOUT,
              });

              setIsEmailAvailable(isAvailable);
              setErrors((prev) => ({
                ...prev,
                email: isAvailable ? undefined : "Email is already registered",
              }));
            } catch (error: unknown) {
              setErrors((prev) => ({
                ...prev,
                email: "Could not validate email availability",
              }));

              // TODO (1)
              //
              // if (error instanceof APIError) {
              //   toast.error(error.message);
              // } else {
              //   toast.error("Internal server error");
              // }
            }
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
    (field: keyof SignUpForm, value: string) => {
      if (field === "email") setIsEmailAvailable(false);

      setForm((prev) => {
        const nextForm: SignUpForm = { ...prev, [field]: value };

        validate(
          field,
          value,
          nextForm.password,
          queryClient,
          setErrors,
          setIsEmailAvailable,
        );

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
  const setConfirmPassword = useCallback(
    (password: string) => update("confirmPassword", password),
    [update],
  );

  const handleSignUp = () => {
    const res: ZodSafeParseResult<SignUpForm> = signUpSchema.safeParse(form);

    if (res.success && form.password !== form.confirmPassword) {
      setErrors((prev) => ({
        ...prev,
        confirmPassword: "Passwords do not match",
      }));
      return;
    } else if (res.success) {
      setErrors({});

      queryClient.invalidateQueries({
        queryKey: ["email-availability", form.email],
      });

      signUp(form, {
        onSuccess: () => {
          setLocalRememberMe(true);
        },
      });
    } else {
      const newErrors: SignUpFormErrors = {};

      res.error.issues.forEach((issue) => {
        const key: PropertyKey = issue.path[0];
        if (key) {
          if (key === "email") setIsEmailAvailable(false);

          newErrors[key as keyof SignUpForm] = issue.message;
        }
      });

      setErrors(newErrors);
    }
  };

  return {
    form,
    errors,
    setEmail,
    isEmailAvailable,
    setPassword,
    setConfirmPassword,
    handleSignUp,
    isSigningUp: isPending,
  };
};

export default useSignUpForm;
