const ROUTES = {
  HOME: "/",
  AUTH: {
    FORGOT_PASSWORD: "/auth/forgot-password",
    RESET_PASSWORD: "/auth/reset-password",
    SIGN_IN: "/auth/sign-in",
    SIGN_UP: "/auth/sign-up",
  },
  LEGAL: {
    PRIVACY_POLICY: "/legal/privacy-policy",
    TOS: "/legal/terms-of-service",
  },
} as const;

export default ROUTES;
