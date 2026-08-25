import { z } from "zod";

const forgotPasswordSchema = z.object({
  email: z.email("Email is invalid").trim(),
});

type ForgotPasswordForm = z.infer<typeof forgotPasswordSchema>;

export { forgotPasswordSchema, type ForgotPasswordForm };
