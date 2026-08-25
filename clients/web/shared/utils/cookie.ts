import Cookies from "js-cookie";

export const setCookieAccessToken = (
  token: string,
  attributes: Cookies.CookieAttributes,
): void => {
  Cookies.set("access_token", token, attributes);
};
