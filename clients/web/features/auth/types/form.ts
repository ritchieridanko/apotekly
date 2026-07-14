import type { SignInForm, SignUpForm } from "@/features/auth/schemas";

type SignInFormErrors = Partial<Record<keyof SignInForm, string>>;
type SignUpFormErrors = Partial<Record<keyof SignUpForm, string>>;

export type { SignInFormErrors, SignUpFormErrors };
