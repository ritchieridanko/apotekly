import { create } from "zustand";

interface Auth {
  email: string;
  role: string;
  isEmailVerified: boolean;
}

interface AuthState {
  auth: Auth | null;
  setAuth: (auth: Auth | null) => void;
  accessToken: string | null;
  setAccessToken: (token: string | null) => void;
  clearAuth: () => void;
}

const useAuthStore = create<AuthState>((set) => ({
  auth: null,
  setAuth: (auth) => set({ auth: auth }),
  accessToken: null,
  setAccessToken: (token) => set({ accessToken: token }),
  clearAuth: () =>
    set({
      auth: null,
      accessToken: null,
    }),
}));

export { useAuthStore };
