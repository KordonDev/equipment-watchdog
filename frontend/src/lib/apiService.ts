import { routes } from "./routes";
import { BASE_URL } from "./constants";
import { goto } from '$app/navigation';
import { showError } from '$lib/services/notification.svelte';


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
				return res.json().then((data: any) => {
					if (data.redirect === "login") {
						return goto(routes.login, { replaceState: true });
					}
					if (data.redirect === "not-approved") {
						return goto(routes.notApproved, { replaceState: true });
					}
				});
			}
			showError('Ein Fehler ist aufgetreten');
			throw res;
		});
}
