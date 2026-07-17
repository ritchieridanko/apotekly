import type {
  ForgotPasswordForm,
  ResetPasswordForm,
  SignInForm,
  SignUpForm,
} from "@/features/auth/schemas";

type ForgotPasswordFormErrors = Partial<
  Record<keyof ForgotPasswordForm, string>
>;
type ResetPasswordFormErrors = Partial<Record<keyof ResetPasswordForm, string>>;
type SignInFormErrors = Partial<Record<keyof SignInForm, string>>;
type SignUpFormErrors = Partial<Record<keyof SignUpForm, string>>;

export type {
  ForgotPasswordFormErrors,
  ResetPasswordFormErrors,
  SignInFormErrors,
  SignUpFormErrors,
};
