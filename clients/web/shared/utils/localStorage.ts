const REMEMBER_ME_KEY: string = "apotekly_remember_me";

export const getLocalRememberMe = (): boolean =>
  localStorage.getItem(REMEMBER_ME_KEY) === "true";

export const setLocalRememberMe = (value: boolean): void =>
  localStorage.setItem(REMEMBER_ME_KEY, String(value));
