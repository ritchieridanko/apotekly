const ROUTES = {
  HOME: "/",
  AUTH: {
    SIGN_IN: "/auth/sign-in",
    SIGN_UP: "/auth/sign-up",
    FORGOT_PASSWORD: "/auth/forgot-password",
  },
  LEGAL: {
    TOS: "/legal/terms-of-service",
    PRIVACY_POLICY: "/legal/privacy-policy",
  },
} as const;

export default ROUTES;
