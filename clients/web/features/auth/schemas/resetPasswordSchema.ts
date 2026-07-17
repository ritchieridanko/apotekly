import { z } from "zod";

import {
  PASSWORD_MAX_LENGTH,
  PASSWORD_MIN_LENGTH,
} from "@/features/auth/constants";

const resetPasswordSchema = z.object({
  password: z
    .string()
    .trim()
    .min(
      PASSWORD_MIN_LENGTH,
      `Password must be at least ${PASSWORD_MIN_LENGTH} characters`,
    )
    .max(
      PASSWORD_MAX_LENGTH,
      `Password must not exceed ${PASSWORD_MAX_LENGTH} characters`,
    )
    .regex(/[a-z]/, "Password must include at least one lowercase letter")
    .regex(/[A-Z]/, "Password must include at least one uppercase letter")
    .regex(/[0-9]/, "Password must include at least one number")
    .regex(
      /[!@#$%^&*()_+\-={};:'"\\|,.<>/?]/,
      "Password must include at least one special character",
    ),
  confirmPassword: z.string().trim().min(1, "Password is invalid"),
});

type ResetPasswordForm = z.infer<typeof resetPasswordSchema>;

export { resetPasswordSchema, type ResetPasswordForm };
