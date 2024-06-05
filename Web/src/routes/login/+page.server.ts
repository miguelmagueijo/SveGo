import type { PageServerLoad } from "./$types";
import { redirect } from "@sveltejs/kit";

export const prerender = false;

export const load: PageServerLoad = ({ cookies }) => {
	const accessToken = cookies.get("token");
	const refreshToken = cookies.get("refreshToken");

	if (accessToken || refreshToken) {
		return redirect(302, "/");
	}
};
