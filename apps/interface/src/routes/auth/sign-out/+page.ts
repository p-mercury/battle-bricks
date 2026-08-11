import type { PageLoad } from "./$types";
import { redirect } from "@sveltejs/kit";
import { External } from "@battle-bricks/contracts/identity/v1/service_pb";

export const load: PageLoad = async ({ url }) => {
	const redirectUrl = new URL(
		`/api/${External.typeName}/auth/signout`,
		url.origin,
	);
	redirectUrl.searchParams.set("redirect", url.origin);
	redirect(302, redirectUrl);
};
