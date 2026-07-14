import type {
  ForgotPasswordForm,
  SignInForm,
  SignUpForm,
} from "@/features/auth/schemas";

type ForgotPasswordFormErrors = Partial<
  Record<keyof ForgotPasswordForm, string>
>;
type SignInFormErrors = Partial<Record<keyof SignInForm, string>>;
type SignUpFormErrors = Partial<Record<keyof SignUpForm, string>>;

export type { ForgotPasswordFormErrors, SignInFormErrors, SignUpFormErrors };
