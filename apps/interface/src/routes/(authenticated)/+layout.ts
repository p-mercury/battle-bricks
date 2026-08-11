import type { LayoutLoad } from "./$types";
import { loadLocale } from "wuchale/load-utils";
import { getClient } from "$lib/clients/univeral-client.svelte";

export const load: LayoutLoad = async (event) => {
	loadLocale(event.data.locale);
	const client = getClient(event);
	return { policyUser: (await client.policy.getUser({})).user };
};
