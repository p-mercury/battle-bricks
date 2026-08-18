import type { LayoutLoad } from "./$types";
import { loadLocale } from "wuchale/load-utils";
import { getClient } from "$lib/clients/univeral-client.svelte";

export const load: LayoutLoad = async (event) => {
	loadLocale(event.data.locale);
	const client = getClient(event);

	const [{ user }, { units }] = await Promise.all([
		client.policy.getUser({}),
		client.catalogue.listUnits({}),
	]);

	return {
		policyUser: user,
		units: Object.fromEntries(units.map((i) => [i.id, i])),
	};
};
