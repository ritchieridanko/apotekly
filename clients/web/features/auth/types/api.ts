interface AccessTokenAPIResponse {
  token: string;
  expires_in_seconds: number;
}

interface AuthAPIResponse {
  email: string;
  role: string;
  is_email_verified: boolean;
}

interface CheckEmailAvailabilityAPIResponse {
  is_available: boolean;
}

interface ResetPasswordAPIRequest {
  email: string;
}

interface ResetPasswordAPIResponse {
  email: string;
}

interface RotateAuthTokenAPIResponse {
  access_token?: AccessTokenAPIResponse;
}

interface SignInAPIRequest {
  email: string;
  password: string;
}

interface SignInAPIResponse {
  auth?: AuthAPIResponse;
  access_token?: AccessTokenAPIResponse;
}

interface SignUpAPIRequest {
  email: string;
  password: string;
}

interface SignUpAPIResponse {
  auth?: AuthAPIResponse;
  access_token?: AccessTokenAPIResponse;
}

export type {
  CheckEmailAvailabilityAPIResponse,
  ResetPasswordAPIRequest,
  ResetPasswordAPIResponse,
  RotateAuthTokenAPIResponse,
  SignInAPIRequest,
  SignInAPIResponse,
  SignUpAPIRequest,
  SignUpAPIResponse,
};
