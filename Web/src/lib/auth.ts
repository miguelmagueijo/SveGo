import { writable } from "svelte/store";

export const isAuthenticated = writable(false);

export async function requestLogout() {
	try {
		const res = await fetch("http://localhost:4555/v1/logout", {
			credentials: "include",
		});

		if (res.redirected) {
			window.location.href = "/signin";
		}
	} catch (e) {
		console.error(e);
	}
}
