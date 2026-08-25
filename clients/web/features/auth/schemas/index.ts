import {
  forgotPasswordSchema,
  type ForgotPasswordForm,
} from "./forgotPasswordSchema";
import {
  resetPasswordSchema,
  type ResetPasswordForm,
} from "./resetPasswordSchema";
import { signInSchema, type SignInForm } from "./signInSchema";
import { signUpSchema, type SignUpForm } from "./signUpSchema";

export type { ForgotPasswordForm, ResetPasswordForm, SignInForm, SignUpForm };

export {
  forgotPasswordSchema,
  resetPasswordSchema,
  signInSchema,
  signUpSchema,
};
