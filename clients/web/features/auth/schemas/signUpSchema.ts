import { z } from "zod";

const PASSWORD_MIN_LENGTH: number = 8;
const PASSWORD_MAX_LENGTH: number = 50;

const signUpSchema = z.object({
  email: z.email("Email is invalid").trim(),
  password: z
    .string()
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
    )
    .trim(),
  confirmPassword: z.string().trim(),
});

type SignUpForm = z.infer<typeof signUpSchema>;

export {
  PASSWORD_MIN_LENGTH,
  PASSWORD_MAX_LENGTH,
  signUpSchema,
  type SignUpForm,
};
