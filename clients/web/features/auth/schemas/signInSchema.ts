import { z } from "zod";

const signInSchema = z.object({
  email: z.email("Email is invalid").trim(),
  password: z.string().trim().min(1, "Password is invalid"),
});

type SignInForm = z.infer<typeof signInSchema>;

export { signInSchema, type SignInForm };
