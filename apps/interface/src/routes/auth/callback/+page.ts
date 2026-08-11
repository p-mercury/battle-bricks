import { redirect } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { External } from "@battle-bricks/contracts/identity/v1/service_pb";
import { browser } from "$app/environment";

export const load: PageLoad = async ({ url, fetch }) => {
	if (browser) {
		const code = url.searchParams.get("code");
		if (!code) {
			redirect(302, url.searchParams.get("state") || "/");
		}

		const codeUrl = new URL(`/api/${External.typeName}/auth/code`, url.origin);
		codeUrl.searchParams.set("code", code);
		codeUrl.searchParams.set(
			"redirectUri",
			new URL("/auth/callback", url.origin).href,
		);
		await fetch(codeUrl);

		redirect(302, url.searchParams.get("state") || "/");
	}
};
