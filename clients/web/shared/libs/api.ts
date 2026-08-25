import { useAuthStore } from "@/features/auth/stores";
import { RotateAuthTokenAPIResponse } from "@/features/auth/types";
import {
  API_BASE_URL,
  APP_ENV,
  getLocalRememberMe,
  setCookieAccessToken,
} from "@/shared/utils";

class APIError extends Error {
  constructor(
    public readonly status: number,
    message: string,
  ) {
    super(message);
    this.name = "APIError";
    Object.setPrototypeOf(this, APIError.prototype);
  }
}

type APIOptions = Omit<RequestInit, "credentials" | "body"> & {
  body?: Record<string, any> | any[] | BodyInit | null;
  requiresAuth?: boolean;
};

const fetcher = async (
  path: string,
  options: APIOptions,
): Promise<Response> => {
  const { body } = options;
  const headers: Headers = new Headers(options.headers);
  const isFormData: boolean = body instanceof FormData;

  if (!isFormData && body !== undefined && body !== null) {
    headers.set("Content-Type", "application/json");
  }

  const requiresAuth: boolean = options.requiresAuth ?? false;
  if (requiresAuth) {
    const accessToken = useAuthStore.getState().accessToken;
    if (accessToken) {
      headers.set("Authorization", `Bearer ${accessToken}`);
    }
  }

  return fetch(`${API_BASE_URL}${path}`, {
    ...options,
    credentials: "include",
    headers: headers,
    body:
      body === undefined || body === null
        ? undefined
        : isFormData ||
            typeof body === "string" ||
            body instanceof Blob ||
            body instanceof URLSearchParams ||
            body instanceof ArrayBuffer ||
            body instanceof ReadableStream
          ? (body as BodyInit)
          : JSON.stringify(body),
  });
};

const parse = async <R>(res: Response): Promise<APIResponse<R>> => {
  const text: string = await res.text();
  return text ? (JSON.parse(text) as APIResponse<R>) : ({} as APIResponse<R>);
};

let isRefreshing: Promise<APIResponse<RotateAuthTokenAPIResponse>> | null =
  null;

const api = async <T>(path: string, options: APIOptions = {}): Promise<T> => {
  let res: Response = await fetcher(path, options);
  if (res.status === 204) return undefined as T;

  let payload: T & APIResponse<any> = (await parse(res)) as T &
    APIResponse<any>;

  if (!res.ok) {
    if (res.status === 401 && payload.message === "Unauthenticated") {
      const { setAccessToken, clearAuth } = useAuthStore.getState();
      const rememberMe: boolean = getLocalRememberMe();

      try {
        if (!isRefreshing) {
          isRefreshing = fetcher("/auth/refresh", {
            method: "POST",
            body: { remember_me: rememberMe },
          }).then((res) => parse<RotateAuthTokenAPIResponse>(res));
        }

        const refreshPayload = await isRefreshing;
        if (
          refreshPayload.status >= 400 ||
          !refreshPayload.data?.access_token
        ) {
          clearAuth();
          throw new APIError(refreshPayload.status, refreshPayload.message);
        }

        const token: string = refreshPayload.data.access_token.token;
        setAccessToken(token);

        const attributes: Cookies.CookieAttributes = {
          secure: APP_ENV === "prod",
          sameSite: "Strict",
        };
        if (rememberMe) {
          const seconds: number =
            refreshPayload.data.access_token.expires_in_seconds;
          attributes.expires = new Date(Date.now() + seconds * 1000);
        }

        setCookieAccessToken(token, attributes);

        // Retry the original request
        res = await fetcher(path, options);
        payload = (await parse(res)) as T & APIResponse<any>;
      } catch (error: unknown) {
        clearAuth();

        throw error instanceof APIError
          ? error
          : new APIError(401, "Unauthenticated");
      } finally {
        isRefreshing = null;
      }
    }
  }
  if (!res.ok) {
    throw new APIError(
      payload.status ?? res.status,
      payload.message ?? res.statusText,
    );
  }

  return payload as T;
};

export { APIError, api };
