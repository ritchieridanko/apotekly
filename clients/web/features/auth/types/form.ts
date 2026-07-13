import type { SignUpForm } from "@/features/auth/schemas";

type SignUpFormErrors = Partial<Record<keyof SignUpForm, string>>;

export type { SignUpFormErrors };
