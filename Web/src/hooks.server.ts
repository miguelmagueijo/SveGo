import { type Handle, redirect } from "@sveltejs/kit";

export const handle: Handle = async ({ event, resolve }) => {
	event.locals.isAuthenticated =
		event.cookies.get("token") !== undefined || event.cookies.get("refreshToken") !== undefined;

	if (event.url.pathname.startsWith("/login") && event.locals.isAuthenticated) {
		redirect(302, "/");
	}

	if (event.url.pathname === "/" && !event.locals.isAuthenticated) {
		redirect(302, "/signin");
	}

	return resolve(event);
};
