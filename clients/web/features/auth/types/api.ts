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

interface RotateAuthTokenAPIResponse {
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
  RotateAuthTokenAPIResponse,
  SignUpAPIRequest,
  SignUpAPIResponse,
};
