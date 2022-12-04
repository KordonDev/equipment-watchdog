import { routes } from "../routes";
import { BASE_URL } from "../constants";
import { replace } from "svelte-spa-router";

export function fetchApi(url: string, headers?: RequestInit) {
  return fetch(`${BASE_URL}${url}`, {
    credentials: "include",
    ...headers,
  })
    .then((res) => {
      if (!res.ok) {
        throw res;
      }
      return res;
    })
    .then((res) => {
      const contentType = res.headers.get("content-type");
      if (contentType && contentType.indexOf("application/json") !== -1) {
        return res.json();
      }
      return res;
    })
    .catch((res) => {
      const contentType = res.headers.get("content-type");
      if (contentType && contentType.indexOf("application/json") !== -1) {
        return res.json().then((data) => {
          if (data.redirect === "login") {
            replace(routes.Login.link);
          }
          if (data.redirect === "not-approved") {
            replace(routes.NotApproved.link);
          }
        });
      }
      throw res;
    });
}
